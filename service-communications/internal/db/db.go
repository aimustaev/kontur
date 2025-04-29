package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"

	"github.com/aimustaev/service-communications/internal/adapter"
)

// DB represents a database connection
type DB struct {
	conn   *pgx.Conn
	logger *logrus.Logger
}

// NewDB creates a new database connection
func NewDB(host string, port int, user, password, dbname string, logger *logrus.Logger) (*DB, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, password, host, port, dbname)
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return &DB{conn: conn, logger: logger}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.conn.Close(context.Background())
}

// SaveEmail implements adapter.Database
func (db *DB) SaveEmail(ctx context.Context, msg adapter.Message) error {
	query := `
		INSERT INTO emails (id, from_address, to_address, subject, body, tags, channel)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := db.conn.Exec(ctx, query, msg.ID, msg.From, msg.To, msg.Subject, msg.Body, msg.Tags, msg.Channel)
	if err != nil {
		return fmt.Errorf("failed to save email: %w", err)
	}
	db.logger.Infof("Saved message %s to database (channel: %s)", msg.ID, msg.Channel)
	return nil
}
