package observability

import (
	"log/slog"
	"os"
	"strings"
)

type Config struct {
	Level   string
	Format  string
	Service string
}

func NewLogger(cfg Config) *slog.Logger {
	var level slog.Level

	switch strings.ToLower(cfg.Level) {
	case "debug":
		level = slog.LevelDebug
	case "warn", "warning":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	options := &slog.HandlerOptions{
		Level: level,
	}

	service := cfg.Service
	if service == "" {
		service = "delta"
	}

	if strings.ToLower(cfg.Format) == "text" {
		return slog.New(
			slog.NewTextHandler(os.Stdout, options),
		).With("service", service)
	}

	return slog.New(
		slog.NewJSONHandler(os.Stdout, options),
	).With("service", service)

}
