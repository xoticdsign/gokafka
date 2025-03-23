package messaging

import (
	"log/slog"

	"gokafka/internal/domain"
	"gokafka/internal/messaging/kafka/consumer"
	"gokafka/internal/messaging/kafka/producer"
)

type App struct {
	Consumer Consumerer
	Producer Producerer
}

type Consumerer interface {
	Run() error
	Shutdown() error
}

type Producerer interface {
	ProduceNotification(notification domain.EventNotification) error
	Shutdown() error
}

func New(cfg domain.MessagingConfig, services *domain.Services, log *slog.Logger) (*App, error) {
	consumer, err := consumer.New(cfg, services, log)
	if err != nil {
		return &App{}, err
	}

	producer, err := producer.New(cfg, log)
	if err != nil {
		return &App{}, err
	}

	return &App{
		Consumer: consumer,
		Producer: producer,
	}, nil
}
