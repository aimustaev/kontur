package repository

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/aimustaev/service-tickets/internal/model"
)

var (
	ErrTicketNotFound = errors.New("ticket not found")
)

// MemoryTicketRepository implements TicketRepository interface with in-memory storage
type MemoryTicketRepository struct {
	tickets map[string]*model.Ticket
	mu      sync.RWMutex
}

// NewMemoryTicketRepository creates a new instance of MemoryTicketRepository
func NewMemoryTicketRepository() *MemoryTicketRepository {
	return &MemoryTicketRepository{
		tickets: make(map[string]*model.Ticket),
	}
}

// Create implements TicketRepository.Create
func (r *MemoryTicketRepository) Create(ctx context.Context, ticket *model.Ticket) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	ticket.CreatedAt = now
	ticket.UpdatedAt = now

	r.tickets[ticket.ID] = ticket
	return nil
}

// GetByID implements TicketRepository.GetByID
func (r *MemoryTicketRepository) GetByID(ctx context.Context, id string) (*model.Ticket, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ticket, exists := r.tickets[id]
	if !exists {
		return nil, ErrTicketNotFound
	}

	return ticket, nil
}

// Update implements TicketRepository.Update
func (r *MemoryTicketRepository) Update(ctx context.Context, ticket *model.Ticket) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tickets[ticket.ID]; !exists {
		return ErrTicketNotFound
	}

	ticket.UpdatedAt = time.Now()
	r.tickets[ticket.ID] = ticket
	return nil
}

// Delete implements TicketRepository.Delete
func (r *MemoryTicketRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tickets[id]; !exists {
		return ErrTicketNotFound
	}

	delete(r.tickets, id)
	return nil
}
