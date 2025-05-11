package repository

import (
	"context"

	"github.com/aimustaev/service-tickets/internal/model"
)

// TicketRepository defines the interface for ticket storage operations
type TicketRepository interface {
	// Create creates a new ticket in the storage
	Create(ctx context.Context, ticket *model.Ticket) error

	// GetByID retrieves a ticket by its ID
	GetByID(ctx context.Context, id string) (*model.Ticket, error)

	// Update updates an existing ticket
	Update(ctx context.Context, ticket *model.Ticket) error

	// Delete removes a ticket by its ID
	Delete(ctx context.Context, id string) error
}
