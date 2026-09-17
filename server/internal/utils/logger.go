package utils

import (
	"log/slog"
	"os"
)

var Log = Logger()

func Logger() *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, nil)

	return slog.New(handler)
}
