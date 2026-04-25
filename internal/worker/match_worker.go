package worker

import (
	"context"
	"log"
	"temukan-api/internal/model"
	"temukan-api/internal/repository"
	"temukan-api/internal/service"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

// MatchJob adalah payload yang dikirim ke worker ketika ada report baru/update.
type MatchJob struct {
	ReportID uuid.UUID
}

// MatchWorker mendengarkan channel jobs dan menjalankan proses matching di background.
type MatchWorker struct {
	jobs         chan MatchJob
	reportRepo   repository.ReportRepository
	matchRepo    repository.MatchRepository
	notifRepo    repository.NotificationRepository
	userRepo     repository.UserRepository
	emailService *service.EmailService // ← Resend
	workerCount  int
}

// NewMatchWorker membuat instance MatchWorker baru.
func NewMatchWorker(
	reportRepo repository.ReportRepository,
	matchRepo repository.MatchRepository,
	notifRepo repository.NotificationRepository,
	userRepo repository.UserRepository,
	emailService *service.EmailService,
	workerCount int,
) *MatchWorker {
	if workerCount < 1 {
		workerCount = 3
	}
	return &MatchWorker{
		jobs:         make(chan MatchJob, 100),
		reportRepo:   reportRepo,
		matchRepo:    matchRepo,
		notifRepo:    notifRepo,
		userRepo:     userRepo,
		emailService: emailService,
		workerCount:  workerCount,
	}
}

// Enqueue mengirim job ke channel (non-blocking).
// Dipanggil dari usecase setelah Create/Update report.
func (w *MatchWorker) Enqueue(reportID uuid.UUID) {
	select {
	case w.jobs <- MatchJob{ReportID: reportID}:
		log.Printf("[MatchWorker] job enqueued for report %s", reportID)
	default:
		log.Printf("[MatchWorker] WARNING: job queue full, dropping job for report %s", reportID)
	}
}

// Start menjalankan worker pool. Blok sampai ctx dibatalkan.
// Panggil dengan `go worker.Start(ctx)` di main.go.
func (w *MatchWorker) Start(ctx context.Context) {
	eg, ctx := errgroup.WithContext(ctx)

	for i := 0; i < w.workerCount; i++ {
		workerID := i
		eg.Go(func() error {
			log.Printf("[MatchWorker] worker #%d started", workerID)
			for {
				select {
				case <-ctx.Done():
					log.Printf("[MatchWorker] worker #%d stopped", workerID)
					return nil
				case job, ok := <-w.jobs:
					if !ok {
						return nil
					}
					w.processJob(ctx, job)
				}
			}
		})
	}

	if err := eg.Wait(); err != nil {
		log.Printf("[MatchWorker] worker error: %v", err)
	}
}

// processJob adalah inti logika matching untuk satu report.
func (w *MatchWorker) processJob(ctx context.Context, job MatchJob) {
	log.Printf("[MatchWorker] processing report %s", job.ReportID)

	report, err := w.reportRepo.FindByID(ctx, job.ReportID)
	if err != nil {
		log.Printf("[MatchWorker] report %s not found: %v", job.ReportID, err)
		return
	}

	if report.Status != model.ReportStatusActive {
		log.Printf("[MatchWorker] report %s is not active, skipping", job.ReportID)
		return
	}

	oppositeType := model.ReportTypeFound
	if report.Type == model.ReportTypeFound {
		oppositeType = model.ReportTypeMissing
	}

	candidates, err := w.reportRepo.FindActiveByType(ctx, oppositeType)
	if err != nil {
		log.Printf("[MatchWorker] failed to fetch candidates: %v", err)
		return
	}

	log.Printf("[MatchWorker] comparing report %s against %d candidates", job.ReportID, len(candidates))

	for _, candidate := range candidates {
		w.tryMatch(ctx, *report, candidate)
	}
}

// tryMatch menghitung skor dan, jika lolos threshold, simpan match + kirim notifikasi.
func (w *MatchWorker) tryMatch(ctx context.Context, report, candidate model.Report) {
	var foundReport, missingReport model.Report
	if report.Type == model.ReportTypeFound {
		foundReport = report
		missingReport = candidate
	} else {
		foundReport = candidate
		missingReport = report
	}

	// Skip duplikat
	exists, err := w.matchRepo.ExistsByReportPair(ctx, foundReport.ID, missingReport.ID)
	if err != nil {
		log.Printf("[MatchWorker] ExistsByReportPair error: %v", err)
		return
	}
	if exists {
		return
	}

	score := service.ScoreReports(foundReport, missingReport)
	log.Printf("[MatchWorker] score found=%s missing=%s => %d", foundReport.ID, missingReport.ID, score)

	if score < service.MinMatchScore {
		return
	}

	// Simpan match ke DB
	savedMatch, err := w.matchRepo.Create(ctx, &model.Match{
		FoundReportID:   foundReport.ID,
		MissingReportID: missingReport.ID,
		Score:           score,
		Notified:        false,
	})
	if err != nil {
		log.Printf("[MatchWorker] failed to save match: %v", err)
		return
	}

	log.Printf("[MatchWorker] match saved: %s (score=%d)", savedMatch.ID, score)

	// Kirim notifikasi DB + email ke kedua pihak
	w.sendNotifications(ctx, savedMatch, foundReport, missingReport)
}

// sendNotifications menyimpan notifikasi ke DB dan mengirim email via Resend.
// Email dikirim sebagai goroutine terpisah (best-effort, tidak block worker).
func (w *MatchWorker) sendNotifications(
	ctx context.Context,
	match *model.Match,
	foundReport, missingReport model.Report,
) {
	matchID := match.ID

	type recipient struct {
		report  model.Report
		message string
		role    string
	}

	recipients := []recipient{
		{
			report:  foundReport,
			message: formatFinderMessage(match.Score, missingReport),
			role:    "finder",
		},
		{
			report:  missingReport,
			message: formatSeekerMessage(match.Score, foundReport),
			role:    "seeker",
		},
	}

	for _, r := range recipients {
		reportID := r.report.ID

		// 1. Simpan ke DB — selalu dilakukan, tidak bergantung email
		notif := &model.Notification{
			UserID:   r.report.ReporterID,
			Message:  r.message,
			IsRead:   false,
			ReportID: &reportID,
			MatchID:  &matchID,
		}
		if err := w.notifRepo.Create(ctx, notif); err != nil {
			log.Printf("[MatchWorker] failed to create DB notification for user %s: %v",
				r.report.ReporterID, err)
		}

		// 2. Kirim email via Resend — best-effort di goroutine terpisah
		userID := r.report.ReporterID
		role := r.role
		go w.sendEmail(ctx, match, foundReport, missingReport, userID, role)
	}

	// Tandai match sudah dinotifikasi
	if err := w.matchRepo.MarkNotified(ctx, matchID); err != nil {
		log.Printf("[MatchWorker] failed to mark match as notified: %v", err)
	}
}

// sendEmail fetch user lalu kirim email. Dijalankan sebagai goroutine terpisah.
func (w *MatchWorker) sendEmail(
	ctx context.Context,
	match *model.Match,
	foundReport, missingReport model.Report,
	userID uuid.UUID,
	role string,
) {
	if w.emailService == nil {
		log.Printf("[MatchWorker] sendEmail: emailService is nil, skipping for user %s", userID)
		return
	}

	user, err := w.userRepo.FindByID(ctx, userID)
	if err != nil {
		log.Printf("[MatchWorker] sendEmail: user %s not found: %v", userID, err)
		return
	}

	payload := service.MatchEmailPayload{
		Match:         match,
		FoundReport:   foundReport,
		MissingReport: missingReport,
		RecipientUser: *user,
		Role:          role,
	}

	if err := w.emailService.SendMatchNotification(ctx, payload); err != nil {
		// Gagal kirim email tidak fatal — notifikasi sudah tersimpan di DB
		log.Printf("[MatchWorker] sendEmail: failed for user %s (%s): %v", userID, role, err)
		return
	}

	log.Printf("[MatchWorker] email sent to %s (%s) for match %s", user.Email, role, match.ID)
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func formatFinderMessage(score int, missingReport model.Report) string {
	name := "seseorang"
	if missingReport.Name != nil && *missingReport.Name != "" {
		name = *missingReport.Name
	}
	return "Laporan penemuan Anda mungkin cocok dengan laporan kehilangan " +
		name + " di " + missingReport.City +
		" (skor: " + itoa(score) + "/100). Cek halaman Matches untuk detail."
}

func formatSeekerMessage(score int, foundReport model.Report) string {
	return "Ada kemungkinan kecocokan untuk laporan Anda! Seseorang ditemukan di " +
		foundReport.City +
		" (skor: " + itoa(score) + "/100). Cek halaman Matches untuk detail."
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := [10]byte{}
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte(n%10) + '0'
		n /= 10
	}
	return string(buf[pos:])
}
