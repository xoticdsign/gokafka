package notificator

import (
	"context"
	"net/http"

	"gokafka/internal/services/notificator/utils/config"
)

type HTTP struct {
	Client Client
}

func New(cfg config.Config) *HTTP {
	return &HTTP{
		Client: Client{
			HTTPClient: http.Client{
				Timeout: cfg.Client.Timeout,
			},
			Handlers: UnimplementedHandlers{},
		},
	}
}

type Client struct {
	HTTPClient http.Client
	Handlers   Handlerer
}

type UnimplementedHandlers struct{}

type Handlerer interface {
	PostNotification(ctx context.Context, value string) error
}

func (u UnimplementedHandlers) PostNotification(ctx context.Context, value string) error {
	return nil
}
