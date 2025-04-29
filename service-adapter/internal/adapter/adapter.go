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
}

// Adapter defines the interface that all adapters must implement
type Adapter interface {
	// Connect establishes a connection to the message source
	Connect(ctx context.Context) error
	// Disconnect closes the connection to the message source
	Disconnect(ctx context.Context) error
	// GetMessages retrieves messages from the source
	GetMessages(ctx context.Context) ([]Message, error)
	// MarkAsProcessed marks a message as processed
	MarkAsProcessed(ctx context.Context, messageID string) error
}

// Database defines the interface for database operations
type Database interface {
	SaveEmail(ctx context.Context, msg Message) error
}
