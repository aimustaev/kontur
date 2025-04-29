package adapter

import "context"

// Message represents a generic message that can be processed by any adapter
type Message struct {
	ID      string
	From    string
	To      string
	Subject string
	Body    string
	Tags    []string
	Channel string // email, telegram
}

// Adapter defines the interface that all adapters must implement
type Adapter interface {
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error
	GetMessages(ctx context.Context) ([]Message, error)
	MarkAsProcessed(ctx context.Context, messageID string) error
}

type Database interface {
	SaveEmail(ctx context.Context, msg Message) error
}
