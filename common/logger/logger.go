package logger

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

func NewLogger() *slog.Logger {
	h := tint.NewHandler(os.Stdout, nil)
	return slog.New(h)
}
