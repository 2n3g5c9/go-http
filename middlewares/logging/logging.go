package logging

import (
	"os"

	"golang.org/x/exp/slog"
)

// Init sets a default structured leveled logger.
func Init(logLevel string) {
	var (
		opts                 = slog.HandlerOptions{Level: parseLevel(logLevel)}
		handler slog.Handler = slog.NewJSONHandler(os.Stdout, &opts)
	)

	// Add git commit to all logs if available.
	if gitCommit != "" {
		handler = handler.WithAttrs([]slog.Attr{slog.String("gitCommit", gitCommit)})
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
}

// parseLevel parses a string into a slog.Level.
func parseLevel(logLevel string) slog.Level {
	switch logLevel {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
