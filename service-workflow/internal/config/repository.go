package config

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// PostgresConfigRepository implements ConfigVersionRepository using PostgreSQL
type PostgresConfigRepository struct {
	db *sqlx.DB
}

// NewPostgresConfigRepository creates a new instance of PostgresConfigRepository
func NewPostgresConfigRepository(db *sqlx.DB) *PostgresConfigRepository {
	return &PostgresConfigRepository{db: db}
}

// GetLatestActive returns the latest active version of a configuration by name
func (r *PostgresConfigRepository) GetLatestActive(name string) (*ConfigVersion, error) {
	query := `
		SELECT id, name, version, content, created_at, updated_at, created_by, is_active
		FROM configs.config_versions
		WHERE name = $1 AND is_active = true
		ORDER BY created_at DESC
		LIMIT 1
	`

	var config ConfigVersion
	err := r.db.Get(&config, query, name)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get latest active config: %w", err)
	}

	return &config, nil
}

// GetByVersion returns a specific version of a configuration
func (r *PostgresConfigRepository) GetByVersion(name, version string) (*ConfigVersion, error) {
	query := `
		SELECT id, name, version, content, created_at, updated_at, created_by, is_active
		FROM configs.config_versions
		WHERE name = $1 AND version = $2
	`

	var config ConfigVersion
	err := r.db.Get(&config, query, name, version)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get config version: %w", err)
	}

	return &config, nil
}

// Create creates a new version of a configuration
func (r *PostgresConfigRepository) Create(config *ConfigVersion) error {
	query := `
		INSERT INTO configs.config_versions (
			id, name, version, content, created_at, updated_at, created_by, is_active
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		)
	`

	if config.ID == uuid.Nil {
		config.ID = uuid.New()
	}
	now := time.Now()
	config.CreatedAt = now
	config.UpdatedAt = now

	_, err := r.db.Exec(query,
		config.ID,
		config.Name,
		config.Version,
		config.Content,
		config.CreatedAt,
		config.UpdatedAt,
		config.CreatedBy,
		config.IsActive,
	)
	if err != nil {
		return fmt.Errorf("failed to create config version: %w", err)
	}

	return nil
}

// Update updates an existing configuration version
func (r *PostgresConfigRepository) Update(config *ConfigVersion) error {
	query := `
		UPDATE configs.config_versions
		SET content = $1, updated_at = $2, is_active = $3
		WHERE name = $4 AND version = $5
	`

	config.UpdatedAt = time.Now()

	result, err := r.db.Exec(query,
		config.Content,
		config.UpdatedAt,
		config.IsActive,
		config.Name,
		config.Version,
	)
	if err != nil {
		return fmt.Errorf("failed to update config version: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("config version not found: %s@%s", config.Name, config.Version)
	}

	return nil
}

// List returns a list of configuration versions matching the filter
func (r *PostgresConfigRepository) List(filter ConfigVersionFilter) ([]*ConfigVersion, error) {
	var conditions []string
	var args []interface{}
	argCount := 1

	if filter.Name != nil {
		conditions = append(conditions, fmt.Sprintf("name = $%d", argCount))
		args = append(args, *filter.Name)
		argCount++
	}

	if filter.Version != nil {
		conditions = append(conditions, fmt.Sprintf("version = $%d", argCount))
		args = append(args, *filter.Version)
		argCount++
	}

	if filter.IsActive != nil {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", argCount))
		args = append(args, *filter.IsActive)
		argCount++
	}

	query := `
		SELECT id, name, version, content, created_at, updated_at, created_by, is_active
		FROM configs.config_versions
	`

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY created_at DESC"

	var configs []*ConfigVersion
	err := r.db.Select(&configs, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list config versions: %w", err)
	}

	return configs, nil
}

// Deactivate deactivates a specific version of a configuration
func (r *PostgresConfigRepository) Deactivate(name, version string) error {
	query := `
		UPDATE configs.config_versions
		SET is_active = false, updated_at = $1
		WHERE name = $2 AND version = $3
	`

	result, err := r.db.Exec(query, time.Now(), name, version)
	if err != nil {
		return fmt.Errorf("failed to deactivate config version: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("config version not found: %s@%s", name, version)
	}

	return nil
}
