package api

import (
	"context"
	"net/http"

	"gokafka/internal/services/api/utils/config"
	"gokafka/internal/services/api/utils/helpers"
)

type HTTP struct {
	Server Server
	Client Client
}

func New(cfg config.Config) *HTTP {
	return &HTTP{
		Server: Server{
			HTTPServer: &http.Server{
				Addr:         helpers.CreateAddress(cfg.Host, cfg.Port),
				ReadTimeout:  cfg.ReadTimeout,
				WriteTimeout: cfg.WriteTimeout,
			},
			API: API{
				Service: &UnimplementedService{},
			},
		},
		Client: Client{
			HTTPClient: http.Client{
				Timeout: cfg.Client.Timeout,
			},
			Handlers: UnimplementedHandlers{},
		},
	}
}

type Server struct {
	HTTPServer *http.Server
	API        API
}

type API struct {
	Service Servicer
}

type UnimplementedService struct{}

type Servicer interface {
	Post(ctx context.Context, value string) error
}

type PostRequest struct {
	Value string `json:"value"`
}

type PostResponse struct {
	Value string `json:"value"`
}

func (u *UnimplementedService) Post(ctx context.Context, value string) error {
	return nil
}

type Client struct {
	HTTPClient http.Client
	Handlers   Handlerer
}

type UnimplementedHandlers struct{}

type Handlerer interface{}
