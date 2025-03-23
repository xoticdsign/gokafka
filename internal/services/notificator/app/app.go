package app

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"gokafka/internal/domain"
	"gokafka/internal/lib/logger"
	"gokafka/internal/messaging"
	"gokafka/internal/services/notificator/app/notificator"
	"gokafka/internal/services/notificator/utils/config"
)

func New(env string) (*App, error) {
	cfg, err := config.LoadConfig(env)
	if err != nil {
		return &App{}, err
	}

	log := logger.Init(env)

	notificator := notificator.New(cfg, log)

	app := &App{
		Notificator: notificator,
		Log:         log,
		Config:      cfg,
	}

	messaging, err := messaging.New(
		domain.MessagingConfig{
			Brokers: cfg.KafkaConfig.Brokers,
			GroupID: cfg.KafkaConfig.GroupID,
			Topics:  cfg.KafkaConfig.Topics,
		},
		&domain.Services{
			Notificator: app.Notificator.Definition.HTTP.Client.Handlers,
		},
		log,
	)
	if err != nil {
		return &App{}, err
	}

	app.Messaging = messaging

	return app, nil
}

type App struct {
	Notificator *notificator.App
	Messaging   *messaging.App
	Log         *slog.Logger
	Config      config.Config
}

func (a *App) Run() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	errChan := make(chan error, 1)

	a.Log.Info(
		"starting notificator",
		slog.Any("config", a.Config),
	)

	go func() {
		a.Log.Info(
			"starting messaging",
			slog.Any("config", a.Config.KafkaConfig),
		)

		err := a.Messaging.Consumer.Run()
		if err != nil {
			errChan <- err
		}
	}()

	select {
	case <-sigChan:
		a.Log.Info(
			"prompt to shutdown",
		)

	case err := <-errChan:
		a.Log.Error(
			"have to shutdown",
			slog.String("reason", err.Error()),
		)
	}

	defer a.shutdown()
}

func (a *App) shutdown() {
	err := a.Messaging.Consumer.Shutdown()
	if err != nil {
		a.Log.Error(
			"messaging shutdown with force",
			slog.String("reason", err.Error()),
		)
	}
	a.Log.Info(
		"messaging shutdown gracefully",
	)

	a.Notificator.Shutdown()

	a.Log.Info(
		"notificator shutdown gracefully",
	)
}

type UnimplementedApp struct{}

func (u *UnimplementedApp) Run() error {
	return nil
}
