package logger

import (
	"log/slog"
	"os"
)

// SetLogger configures the global slog logger based on the application environment.
func SetLogger(appEnv string) {
	var logger *slog.Logger
	if appEnv == "production" {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug, // Define o nível mínimo de log a ser exibido
		}))
	} else {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}
	slog.SetDefault(logger)
}
