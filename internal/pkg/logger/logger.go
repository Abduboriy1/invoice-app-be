package logger

import (
	"io"
	"log/slog"
	"os"
)

type Logger struct {
	*slog.Logger
}

func New(level string) *Logger {
	var lvl slog.Level
	switch level {
	case "debug":
		lvl = slog.LevelDebug
	case "info":
		lvl = slog.LevelInfo
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}

	// Open debug log file
	f, _ := os.OpenFile("debug.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	// Write to BOTH stdout AND file
	multiWriter := io.MultiWriter(os.Stdout, f)

	handler := slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
		Level: lvl,
	})

	return &Logger{
		Logger: slog.New(handler),
	}
}

func (l *Logger) WithField(key string, value any) *Logger {
	return &Logger{
		Logger: l.Logger.With(key, value),
	}
}
