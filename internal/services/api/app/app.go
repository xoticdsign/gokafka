package app

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"gokafka/internal/domain"
	"gokafka/internal/lib/logger"
	"gokafka/internal/messaging"
	"gokafka/internal/services/api/app/api"
	"gokafka/internal/services/api/utils/config"
)

type App struct {
	API       *api.App
	Messaging *messaging.App
	Log       *slog.Logger
	Config    config.Config
}

func New(env string) (*App, error) {
	cfg, err := config.LoadConfig(env)
	if err != nil {
		return &App{}, err
	}

	log := logger.Init(env)

	messaging, err := messaging.New(
		domain.MessagingConfig{
			Brokers: cfg.KafkaConfig.Brokers,
			GroupID: cfg.KafkaConfig.GroupID,
			Topics:  cfg.KafkaConfig.Topics,
		},
		&domain.Services{},
		log,
	)
	if err != nil {
		return &App{}, err
	}

	api, err := api.New(cfg, messaging, log)
	if err != nil {
		return &App{}, err
	}

	return &App{
		API:       api,
		Messaging: messaging,
		Log:       log,
		Config:    cfg,
	}, nil
}

func (a *App) Run() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	errChan := make(chan error, 1)

	go func() {
		a.Log.Info(
			"starting api",
			slog.Any("config", a.Config),
		)

		err := a.API.Run()
		if err != nil {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		a.Log.Error(
			"have to shutdown",
			slog.String("reason", err.Error()),
		)

	case <-sigChan:
		a.Log.Info(
			"prompt to shutdown",
		)
	}

	defer a.shutdown()
}

func (a *App) shutdown() {
	err := a.Messaging.Producer.Shutdown()
	if err != nil {
		a.Log.Error(
			"messaging producer shutdown with force",
			slog.String("reason", err.Error()),
		)
	}
	a.Log.Info(
		"messaging producer shutdown gracefully",
	)

	err = a.API.Shutdown()
	if err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			a.Log.Error(
				"api shutdown timeout exceeded, using force",
				slog.String("reason", err.Error()),
			)

		default:
			a.Log.Error(
				"api shutdown with force",
				slog.String("reason", err.Error()),
			)
		}
	}
	a.Log.Info(
		"api shutdown gracefully",
	)
}

type UnimplementedApp struct{}

func (u *UnimplementedApp) Run() {}
