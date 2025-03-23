package logger

import (
	"log/slog"
	"os"
)

func Init(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case "prod":
		// TODO

	case "dev":
		// TODO

	default:
		log = slog.New(
			slog.NewTextHandler(
				os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelDebug,
				},
			),
		)
	}

	return log
}
