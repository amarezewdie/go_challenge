package customlogger

import (
	"log/slog"
	"os"
)

var LogLevels = map[string]slog.Level{
	"debug": slog.LevelDebug,
	"info":  slog.LevelInfo,
	"warn":  slog.LevelWarn,
	"error": slog.LevelError,
}

func NewLogger(env string, lvl slog.Level, version string) *slog.Logger {
	var logHandler slog.Handler
	switch env {
	case "development":
		logHandler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: lvl,
		})
	default:
		logHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: lvl,
		})
	}

	logger := slog.New(logHandler).With(slog.Group(
		"service_info",
		slog.String("env", env),
		slog.String("version", version),
	),
	)

	return logger
}
