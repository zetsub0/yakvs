package logger

import (
	"log/slog"
	"os"
)

// SetupLogger creates slog.Logger with json type
func SetupLogger(env string) *slog.Logger {

	opts := &slog.HandlerOptions{
		AddSource: true,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == "file" {
				return slog.Attr{}
			}
			return a
		},
	}

	switch env {
	case "dev":
		opts.Level = slog.LevelDebug
	case "local":
		opts.Level = slog.LevelDebug
	case "prod":
		opts.Level = slog.LevelError
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)

	logger := slog.New(handler)
	return logger
}
