package logger

import (
	"log/slog"
	"os"
	"sync"
)

var (
	instance *slog.Logger
	once     sync.Once
)

// Init menginisialisasi global logger (JSON, level dari env LOG_LEVEL).
// Panggil sekali di main sebelum komponen lain digunakan.
func Init() {
	once.Do(func() {
		level := slog.LevelInfo
		if os.Getenv("LOG_LEVEL") == "debug" {
			level = slog.LevelDebug
		}

		handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     level,
			AddSource: true, // tampilkan file:line di setiap log
		})

		instance = slog.New(handler)
		slog.SetDefault(instance)
	})
}

// Get mengembalikan global logger instance.
func Get() *slog.Logger {
	if instance == nil {
		Init()
	}
	return instance
}

// ── Shorthand helpers ─────────────────────────────────────────────────────────

func Info(msg string, args ...any)  { Get().Info(msg, args...) }
func Warn(msg string, args ...any)  { Get().Warn(msg, args...) }
func Error(msg string, args ...any) { Get().Error(msg, args...) }
func Debug(msg string, args ...any) { Get().Debug(msg, args...) }

// Fatal log error lalu exit — gunakan hanya di startup/critical path.
func Fatal(msg string, args ...any) {
	Get().Error(msg, args...)
	os.Exit(1)
}