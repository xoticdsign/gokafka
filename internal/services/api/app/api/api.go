package api

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"gokafka/internal/messaging"
	"gokafka/internal/services/api/app/api/definition"
	"gokafka/internal/services/api/utils/config"
)

func New(cfg config.Config, messaging *messaging.App, log *slog.Logger) (*App, error) {
	definition := definition.New(cfg, messaging, log)

	return &App{
		Definition: definition,
		Log:        log,
	}, nil
}

type App struct {
	Definition *definition.Definition
	Log        *slog.Logger
}

func (a *App) Run() error {
	err := a.Definition.HTTP.Server.HTTPServer.ListenAndServe()
	if err != nil {
		if err != http.ErrServerClosed {
			return err
		}
	}
	return nil
}

func (a *App) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err := a.Definition.HTTP.Server.HTTPServer.Shutdown(ctx)
	if err != nil {
		return err
	}
	return nil
}
