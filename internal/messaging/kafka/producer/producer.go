package producer

import (
	"encoding/json"
	"log/slog"

	"github.com/IBM/sarama"

	"gokafka/internal/domain"
)

func New(cfg domain.MessagingConfig, log *slog.Logger) (*Producer, error) {
	cfgSarama := sarama.NewConfig()

	cfgSarama.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(
		cfg.Brokers,
		cfgSarama,
	)
	if err != nil {
		return &Producer{}, err
	}

	return &Producer{
		Messenger: producer,
		Log:       log,
		Config:    cfg,
	}, nil
}

type Producer struct {
	Messenger sarama.SyncProducer
	Log       *slog.Logger
	Config    domain.MessagingConfig
}

func (p *Producer) ProduceNotification(notification domain.EventNotification) error {
	toSend, _ := json.Marshal(notification)

	_, _, err := p.Messenger.SendMessage(
		&sarama.ProducerMessage{
			Topic: "my-topic",
			Value: sarama.ByteEncoder(toSend),
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (p *Producer) Shutdown() error {
	err := p.Messenger.Close()
	if err != nil {
		return err
	}
	return nil
}

type UnimplementedProducer struct{}

func (p *UnimplementedProducer) ProduceNotification(notification domain.EventNotification) error {
	return nil
}

func (p *UnimplementedProducer) Shutdown() error {
	return nil
}
