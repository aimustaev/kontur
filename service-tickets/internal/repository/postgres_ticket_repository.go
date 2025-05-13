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
			status VARCHAR(50) NOT NULL,
			"user" VARCHAR(255) NOT NULL,
			agent VARCHAR(255),
			problem_id BIGINT,
			vertical_id BIGINT,
			skill_id BIGINT,
			user_group_id BIGINT,
			channel VARCHAR(50) NOT NULL,
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
		INSERT INTO tickets (
			id, status, "user", agent, problem_id, vertical_id, 
			skill_id, user_group_id, channel, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err := r.db.ExecContext(ctx, query,
		ticket.ID,
		ticket.Status,
		ticket.User,
		ticket.Agent,
		ticket.ProblemID,
		ticket.VerticalID,
		ticket.SkillID,
		ticket.UserGroupID,
		ticket.Channel,
		ticket.CreatedAt,
		ticket.UpdatedAt,
	)

	return err
}

// GetByID implements TicketRepository.GetByID
func (r *PostgresTicketRepository) GetByID(ctx context.Context, id string) (*model.Ticket, error) {
	query := `
		SELECT 
			id, status, "user", agent, problem_id, vertical_id,
			skill_id, user_group_id, channel, created_at, updated_at
		FROM tickets
		WHERE id = $1
	`

	ticket := &model.Ticket{}
	var agent sql.NullString
	var problemID, verticalID, skillID, userGroupID sql.NullInt64

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&ticket.ID,
		&ticket.Status,
		&ticket.User,
		&agent,
		&problemID,
		&verticalID,
		&skillID,
		&userGroupID,
		&ticket.Channel,
		&ticket.CreatedAt,
		&ticket.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrTicketNotFound
	}
	if err != nil {
		return nil, err
	}

	// Convert nullable fields to pointers
	if agent.Valid {
		ticket.Agent = &agent.String
	}
	if problemID.Valid {
		val := problemID.Int64
		ticket.ProblemID = &val
	}
	if verticalID.Valid {
		val := verticalID.Int64
		ticket.VerticalID = &val
	}
	if skillID.Valid {
		val := skillID.Int64
		ticket.SkillID = &val
	}
	if userGroupID.Valid {
		val := userGroupID.Int64
		ticket.UserGroupID = &val
	}

	return ticket, nil
}

// Update implements TicketRepository.Update
func (r *PostgresTicketRepository) Update(ctx context.Context, ticket *model.Ticket) error {
	ticket.UpdatedAt = time.Now()

	query := `
		UPDATE tickets
		SET 
			status = $1,
			"user" = $2,
			agent = $3,
			problem_id = $4,
			vertical_id = $5,
			skill_id = $6,
			user_group_id = $7,
			channel = $8,
			updated_at = $9
		WHERE id = $10
	`

	result, err := r.db.ExecContext(ctx, query,
		ticket.Status,
		ticket.User,
		ticket.Agent,
		ticket.ProblemID,
		ticket.VerticalID,
		ticket.SkillID,
		ticket.UserGroupID,
		ticket.Channel,
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

// GetActiveByUser implements TicketRepository.GetActiveByUser
func (r *PostgresTicketRepository) GetActiveByUser(ctx context.Context, user string) ([]*model.Ticket, error) {
	query := `
		SELECT 
			id, status, "user", agent, problem_id, vertical_id,
			skill_id, user_group_id, channel, created_at, updated_at
		FROM tickets
		WHERE "user" = $1 AND status NOT IN ('closed', 'resolved')
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, user)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []*model.Ticket
	for rows.Next() {
		ticket := &model.Ticket{}
		var agent sql.NullString
		var problemID, verticalID, skillID, userGroupID sql.NullInt64

		err := rows.Scan(
			&ticket.ID,
			&ticket.Status,
			&ticket.User,
			&agent,
			&problemID,
			&verticalID,
			&skillID,
			&userGroupID,
			&ticket.Channel,
			&ticket.CreatedAt,
			&ticket.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Convert nullable fields to pointers
		if agent.Valid {
			ticket.Agent = &agent.String
		}
		if problemID.Valid {
			val := problemID.Int64
			ticket.ProblemID = &val
		}
		if verticalID.Valid {
			val := verticalID.Int64
			ticket.VerticalID = &val
		}
		if skillID.Valid {
			val := skillID.Int64
			ticket.SkillID = &val
		}
		if userGroupID.Valid {
			val := userGroupID.Int64
			ticket.UserGroupID = &val
		}

		tickets = append(tickets, ticket)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tickets, nil
}

// GetAll implements TicketRepository.GetAll
func (r *PostgresTicketRepository) GetAll(ctx context.Context) ([]*model.Ticket, error) {
	query := `
		SELECT 
			id, status, "user", agent, problem_id, vertical_id,
			skill_id, user_group_id, channel, created_at, updated_at
		FROM tickets
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []*model.Ticket
	for rows.Next() {
		ticket := &model.Ticket{}
		var agent sql.NullString
		var problemID, verticalID, skillID, userGroupID sql.NullInt64

		err := rows.Scan(
			&ticket.ID,
			&ticket.Status,
			&ticket.User,
			&agent,
			&problemID,
			&verticalID,
			&skillID,
			&userGroupID,
			&ticket.Channel,
			&ticket.CreatedAt,
			&ticket.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Convert nullable fields to pointers
		if agent.Valid {
			ticket.Agent = &agent.String
		}
		if problemID.Valid {
			val := problemID.Int64
			ticket.ProblemID = &val
		}
		if verticalID.Valid {
			val := verticalID.Int64
			ticket.VerticalID = &val
		}
		if skillID.Valid {
			val := skillID.Int64
			ticket.SkillID = &val
		}
		if userGroupID.Valid {
			val := userGroupID.Int64
			ticket.UserGroupID = &val
		}

		tickets = append(tickets, ticket)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tickets, nil
}

// Close closes the database connection
func (r *PostgresTicketRepository) Close() error {
	return r.db.Close()
}
