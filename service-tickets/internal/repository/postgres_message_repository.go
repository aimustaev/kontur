package repository

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/aimustaev/service-tickets/internal/model"
)

// PostgresMessageRepository implements MessageRepository interface with PostgreSQL storage
type PostgresMessageRepository struct {
	db *sql.DB
}

// NewPostgresMessageRepository creates a new instance of PostgresMessageRepository
func NewPostgresMessageRepository(dsn string) (*PostgresMessageRepository, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Create messages table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS messages (
			id VARCHAR(36) PRIMARY KEY,
			ticket_id VARCHAR(36) NOT NULL,
			from_address VARCHAR(255) NOT NULL,
			to_address VARCHAR(255) NOT NULL,
			subject TEXT NOT NULL,
			body TEXT NOT NULL,
			channel VARCHAR(50) NOT NULL,
			created_at TIMESTAMP NOT NULL,
			FOREIGN KEY (ticket_id) REFERENCES tickets(id)
		)
	`)
	if err != nil {
		return nil, err
	}

	return &PostgresMessageRepository{db: db}, nil
}

// Create implements MessageRepository.Create
func (r *PostgresMessageRepository) Create(ctx context.Context, message *model.Message) error {
	query := `
		INSERT INTO messages (
			id, ticket_id, from_address, to_address, 
			subject, body, channel, created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(ctx, query,
		message.ID,
		message.TicketID,
		message.FromAddress,
		message.ToAddress,
		message.Subject,
		message.Body,
		message.Channel,
		message.CreatedAt,
	)

	return err
}

// GetByTicketID implements MessageRepository.GetByTicketID
func (r *PostgresMessageRepository) GetByTicketID(ctx context.Context, ticketID string) ([]*model.Message, error) {
	query := `
		SELECT id, ticket_id, from_address, to_address, 
			subject, body, channel, created_at
		FROM messages
		WHERE ticket_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, ticketID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*model.Message
	for rows.Next() {
		message := &model.Message{}
		err := rows.Scan(
			&message.ID,
			&message.TicketID,
			&message.FromAddress,
			&message.ToAddress,
			&message.Subject,
			&message.Body,
			&message.Channel,
			&message.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

// GetByID implements MessageRepository.GetByID
func (r *PostgresMessageRepository) GetByID(ctx context.Context, id string) (*model.Message, error) {
	query := `
		SELECT id, ticket_id, from_address, to_address, 
			subject, body, channel, created_at
		FROM messages
		WHERE id = $1
	`

	message := &model.Message{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&message.ID,
		&message.TicketID,
		&message.FromAddress,
		&message.ToAddress,
		&message.Subject,
		&message.Body,
		&message.Channel,
		&message.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrMessageNotFound
	}
	if err != nil {
		return nil, err
	}

	return message, nil
}

// Close closes the database connection
func (r *PostgresMessageRepository) Close() error {
	return r.db.Close()
}
