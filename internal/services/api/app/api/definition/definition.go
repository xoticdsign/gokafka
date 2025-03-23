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
	Service *Service
	Log     *slog.Logger
}

func (a *API) Post(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}
	defer r.Body.Close()

	var post api.PostRequest

	err = json.Unmarshal(body, &post)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	err = a.Service.Post(
		context.Background(),
		post.Value,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	resp, _ := json.Marshal(api.PostResponse{
		Value: "post created",
	})

	w.Header().Set("Content-Type", "application/json")
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
