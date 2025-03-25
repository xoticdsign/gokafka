package definition

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"gokafka/internal/domain"
	"gokafka/internal/messaging"
	"gokafka/internal/services/api/http/api"
	"gokafka/internal/services/api/utils/config"
)

type Definition struct {
	HTTP    *api.HTTP
	Service *Service
}

func New(cfg config.Config, messaging *messaging.App, log *slog.Logger) *Definition {
	service := &Service{
		Messaging: messaging,
		Log:       log,
	}

	mux := http.NewServeMux()

	a := api.New(cfg)

	a.Server.HTTPServer.Handler = mux
	a.Server.API.Service = service

	api := &API{
		Service: service,
		Log:     log,
	}

	mux.HandleFunc("/post", api.Post)

	return &Definition{
		HTTP:    a,
		Service: service,
	}
}

type API struct {
	Service api.Servicer
	Log     *slog.Logger
}

func (a *API) Post(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		resp, _ := json.Marshal(api.PostResponse{
			Value: "method not allowed",
		})

		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(resp)

		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		resp, _ := json.Marshal(api.PostResponse{
			Value: "bad request",
		})

		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp)

		return
	}
	defer r.Body.Close()

	var post api.PostRequest

	err = json.Unmarshal(body, &post)
	if err != nil {
		resp, _ := json.Marshal(api.PostResponse{
			Value: "bad request",
		})

		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp)

		return
	}

	err = a.Service.Post(
		context.Background(),
		post.Value,
	)
	if err != nil {
		resp, _ := json.Marshal(api.PostResponse{
			Value: "internal error",
		})

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(resp)

		return
	}

	resp, _ := json.Marshal(api.PostResponse{
		Value: "post created",
	})

	w.Write(resp)
}

type Service struct {
	Messaging *messaging.App
	Log       *slog.Logger
}

func (s *Service) Post(ctx context.Context, value string) error {
	// hypothetical posting happens here

	fmt.Println(value)

	err := s.Messaging.Producer.ProduceNotification(domain.EventNotification{
		Value: value,
	})
	if err != nil {
		s.Log.Error(
			"messaging producer can't produce",
			slog.String("reason", err.Error()),
		)

		return err
	}
	return nil
}
