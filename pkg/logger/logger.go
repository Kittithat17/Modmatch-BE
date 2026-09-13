// Package logger provides a shared structured logger.
package logger

import (
	"log/slog"
	"os"
)

// New returns a JSON slog logger writing to stdout.
func New() *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	return slog.New(handler)
}
