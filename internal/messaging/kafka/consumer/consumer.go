package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/IBM/sarama"

	"gokafka/internal/domain"
)

func New(cfg domain.MessagingConfig, services *domain.Services, log *slog.Logger) (*Consumer, error) {
	cfgSarama := sarama.NewConfig()

	group, err := sarama.NewConsumerGroup(
		cfg.Brokers,
		cfg.GroupID,
		cfgSarama,
	)
	if err != nil {
		return &Consumer{}, err
	}

	return &Consumer{
		Meesenger: group,
		Services:  services,
		Log:       log,
		Config:    cfg,
	}, nil
}

type Consumer struct {
	Meesenger sarama.ConsumerGroup
	Services  *domain.Services
	Log       *slog.Logger
	Config    domain.MessagingConfig
}

func (c *Consumer) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := c.Meesenger.Consume(
		ctx,
		c.Config.Topics,
		Handlers{
			Services: c.Services,
			Log:      c.Log,
		},
	)
	if err != nil {
		if !errors.Is(err, sarama.ErrClosedConsumerGroup) {
			return err
		}
	}
	return nil
}

func (c *Consumer) Shutdown() error {
	err := c.Meesenger.Close()
	if err != nil {
		return err
	}
	return nil
}

type Handlers struct {
	Services *domain.Services
	Log      *slog.Logger
}

func (h Handlers) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h Handlers) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h Handlers) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case msg := <-claim.Messages():
			switch msg.Topic {
			case "my-topic":
				var eventNotification domain.EventNotification

				json.Unmarshal(msg.Value, &eventNotification)

				err := h.Services.Notificator.PostNotification(context.Background(), eventNotification.Value)
				if err != nil {
					h.Log.Error(
						"notificator returned error",
						slog.String("reason", err.Error()),
					)
				}
			}

			session.MarkMessage(msg, "")

		case <-session.Context().Done():
			return nil
		}
	}
}
