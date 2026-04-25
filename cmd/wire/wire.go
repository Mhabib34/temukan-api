//go:build wireinject
// +build wireinject

package wire

import (
	"context"
	"log"
	"os"

	"temukan-api/config"
	"temukan-api/internal/handler"
	"temukan-api/internal/repository"
	"temukan-api/internal/router"
	"temukan-api/internal/service"
	"temukan-api/internal/usecase"
	"temukan-api/internal/worker"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
)

// ── Provider Sets ─────────────────────────────────────────────────────────────

var infraSet = wire.NewSet(
	config.Load,
	config.NewDB,
	config.NewCloudinary,
	provideValidator,
	provideEmailService,
)

// Constructor sudah return interface — tidak perlu wire.Bind
var repositorySet = wire.NewSet(
	repository.NewUserRepository,
	repository.NewReportRepository,
	repository.NewMatchRepository,
	repository.NewNotificationRepository,
)

var workerSet = wire.NewSet(
	provideMatchWorker,
)

// Constructor sudah return interface — tidak perlu wire.Bind
var usecaseSet = wire.NewSet(
	usecase.NewUserUsecase,
	provideReportUsecase,
	usecase.NewMatchUsecase,
	usecase.NewNotificationUsecase,
)

// Constructor return *Impl — perlu wire.Bind ke interface
var handlerSet = wire.NewSet(
	handler.NewUserHandlerImpl,
	wire.Bind(new(handler.UserHandler), new(*handler.UserHandlerImpl)),

	handler.NewReportHandlerImpl,
	wire.Bind(new(handler.ReportHandler), new(*handler.ReportHandlerImpl)),

	handler.NewMatchHandlerImpl,
	wire.Bind(new(handler.MatchHandler), new(*handler.MatchHandlerImpl)),

	handler.NewNotificationHandlerImpl,
	wire.Bind(new(handler.NotificationHandler), new(*handler.NotificationHandlerImpl)),
)

// ── Individual Providers ──────────────────────────────────────────────────────

func provideValidator() *validator.Validate {
	return validator.New()
}

func provideEmailService() *service.EmailService {
	return service.NewEmailService(
		os.Getenv("RESEND_API_KEY"),
		os.Getenv("RESEND_FROM_EMAIL"),
		os.Getenv("APP_URL"),
	)
}

// provideMatchWorker membuat MatchWorker, langsung start di goroutine,
// dan return cleanup func untuk graceful shutdown.
func provideMatchWorker(
	reportRepo repository.ReportRepository,
	matchRepo repository.MatchRepository,
	notifRepo repository.NotificationRepository,
	userRepo repository.UserRepository,
	emailSvc *service.EmailService,
) (*worker.MatchWorker, func(), error) {
	mw := worker.NewMatchWorker(reportRepo, matchRepo, notifRepo, userRepo, emailSvc, 3)

	ctx, cancel := context.WithCancel(context.Background())
	go mw.Start(ctx)
	log.Println("[MatchWorker] started")

	cleanup := func() {
		cancel()
		log.Println("[MatchWorker] stopped")
	}

	return mw, cleanup, nil
}

// provideReportUsecase adalah wrapper karena NewReportUsecase memakai variadic
// parameter untuk matchWorker — Wire tidak bisa inject variadic secara langsung.
func provideReportUsecase(
	repo repository.ReportRepository,
	validate *validator.Validate,
	cld *cloudinary.Cloudinary,
	mw *worker.MatchWorker,
) usecase.ReportUsecase {
	return usecase.NewReportUsecase(repo, validate, cld, mw)
}

// ── Injector ──────────────────────────────────────────────────────────────────

// InitializeApp adalah entry-point Wire.
// Jalankan: go run github.com/google/wire/cmd/wire@latest gen ./cmd/wire/...
func InitializeApp() (*gin.Engine, func(), error) {
	wire.Build(
		infraSet,
		repositorySet,
		workerSet,
		usecaseSet,
		handlerSet,
		router.SetupRouter,
	)
	return nil, nil, nil
}
