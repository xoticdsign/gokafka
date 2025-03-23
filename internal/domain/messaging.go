package domain

import "context"

type MessagingConfig struct {
	Brokers []string
	GroupID string
	Topics  []string
}

type Services struct {
	API         APIHandlerer
	Notificator NotificatorHandlerer
}

type APIHandlerer interface {
	Post(ctx context.Context, value string) error
}

type NotificatorHandlerer interface {
	PostNotification(ctx context.Context, value string) error
}

type EventNotification struct {
	Value string `json:"value"`
}
