package definition

import (
	"context"
	"fmt"
	"log/slog"

	"gokafka/internal/services/notificator/http/notificator"
	"gokafka/internal/services/notificator/utils/config"
)

type Definition struct {
	HTTP *notificator.HTTP
}

func New(cfg config.Config, log *slog.Logger) *Definition {
	n := notificator.New(cfg)

	n.Client.Handlers = Handlers{
		Log: log,
	}

	return &Definition{
		HTTP: n,
	}
}

type Handlers struct {
	Log *slog.Logger
}

func (c Handlers) PostNotification(ctx context.Context, value string) error {
	// hypothetical notifications are sending here

	fmt.Println(value)

	return nil
}
