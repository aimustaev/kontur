package repository

import (
	"context"

	"github.com/aimustaev/service-tickets/internal/model"
)

// MessageRepository defines the interface for message storage operations
type MessageRepository interface {
	// Create creates a new message in the storage
	Create(ctx context.Context, message *model.Message) error

	// GetByTicketID retrieves all messages for a given ticket
	GetByTicketID(ctx context.Context, ticketID string) ([]*model.Message, error)

	// GetByID retrieves a message by its ID
	GetByID(ctx context.Context, id string) (*model.Message, error)
}
