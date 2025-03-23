package notificator

import (
	"log/slog"

	"gokafka/internal/services/notificator/app/notificator/definition"
	"gokafka/internal/services/notificator/utils/config"
)

func New(cfg config.Config, log *slog.Logger) *App {
	definition := definition.New(cfg, log)

	return &App{
		Definition: definition,
		Log:        log,
	}
}

type App struct {
	Definition *definition.Definition
	Log        *slog.Logger
}

func (a *App) Shutdown() {
	a.Definition.HTTP.Client.HTTPClient.CloseIdleConnections()
}
