package utils

import (
	"log/slog"
	"os"
)

func SetupLogger(level string, format string) error {
	l := &slog.LevelVar{}
	l.UnmarshalText([]byte(level))

	opts := slog.HandlerOptions{
		Level: l,
	}

	var h slog.Handler
	h = slog.NewTextHandler(os.Stdout, &opts)
	if format == "json" {
		h = slog.NewJSONHandler(os.Stdout, &opts)
	}

	slog.SetDefault(slog.New(h))
	return nil
}
