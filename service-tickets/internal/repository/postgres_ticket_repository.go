package repository

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"

	"github.com/aimustaev/service-tickets/internal/model"
)

// PostgresTicketRepository implements TicketRepository interface with PostgreSQL storage
type PostgresTicketRepository struct {
	db *sql.DB
}

// NewPostgresTicketRepository creates a new instance of PostgresTicketRepository
func NewPostgresTicketRepository(dsn string) (*PostgresTicketRepository, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Create tickets table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tickets (
			id VARCHAR(36) PRIMARY KEY,
			vertical_id VARCHAR(36) NOT NULL,
			user_id VARCHAR(36) NOT NULL,
			assign VARCHAR(36),
			skill_id VARCHAR(36),
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)
	`)
	if err != nil {
		return nil, err
	}

	return &PostgresTicketRepository{db: db}, nil
}

// Create implements TicketRepository.Create
func (r *PostgresTicketRepository) Create(ctx context.Context, ticket *model.Ticket) error {
	now := time.Now()
	ticket.CreatedAt = now
	ticket.UpdatedAt = now

	query := `
		INSERT INTO tickets (id, vertical_id, user_id, assign, skill_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(ctx, query,
		ticket.ID,
		ticket.VerticalID,
		ticket.UserID,
		ticket.Assign,
		ticket.SkillID,
		ticket.CreatedAt,
		ticket.UpdatedAt,
	)

	return err
}

// GetByID implements TicketRepository.GetByID
func (r *PostgresTicketRepository) GetByID(ctx context.Context, id string) (*model.Ticket, error) {
	query := `
		SELECT id, vertical_id, user_id, assign, skill_id, created_at, updated_at
		FROM tickets
		WHERE id = $1
	`

	ticket := &model.Ticket{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&ticket.ID,
		&ticket.VerticalID,
		&ticket.UserID,
		&ticket.Assign,
		&ticket.SkillID,
		&ticket.CreatedAt,
		&ticket.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrTicketNotFound
	}
	if err != nil {
		return nil, err
	}

	return ticket, nil
}

// Update implements TicketRepository.Update
func (r *PostgresTicketRepository) Update(ctx context.Context, ticket *model.Ticket) error {
	ticket.UpdatedAt = time.Now()

	query := `
		UPDATE tickets
		SET vertical_id = $1, user_id = $2, assign = $3, skill_id = $4, updated_at = $5
		WHERE id = $6
	`

	result, err := r.db.ExecContext(ctx, query,
		ticket.VerticalID,
		ticket.UserID,
		ticket.Assign,
		ticket.SkillID,
		ticket.UpdatedAt,
		ticket.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrTicketNotFound
	}

	return nil
}

// Delete implements TicketRepository.Delete
func (r *PostgresTicketRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM tickets WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrTicketNotFound
	}

	return nil
}

// Close closes the database connection
func (r *PostgresTicketRepository) Close() error {
	return r.db.Close()
}
