package logging

import (
	"log/slog"
	"os"
)

type Logger struct {
	*slog.Logger
}

func (l *Logger) ErrorErr(err error, args ...any) {
	l.Error(err.Error(), args...)
}

func GetLogger(env string, logLevel string) *Logger {

	level := slog.LevelDebug
	switch logLevel {
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}

	if env == "dev" {
		return &Logger{
			slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})),
		}
	}

	return &Logger{
		slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})),
	}
}
