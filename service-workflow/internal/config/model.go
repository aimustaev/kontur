package config

import (
	"time"

	"github.com/google/uuid"
)

// ConfigVersion represents a version of a configuration stored in the database
type ConfigVersion struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Version   string    `json:"version" db:"version"`
	Content   []byte    `json:"content" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy string    `json:"created_by" db:"created_by"`
	IsActive  bool      `json:"is_active" db:"is_active"`
}

// ConfigVersionFilter represents filters for querying config versions
type ConfigVersionFilter struct {
	Name     *string
	Version  *string
	IsActive *bool
}

// ConfigVersionRepository defines the interface for working with config versions
type ConfigVersionRepository interface {
	// GetLatestActive returns the latest active version of a configuration by name
	GetLatestActive(name string) (*ConfigVersion, error)

	// GetByVersion returns a specific version of a configuration
	GetByVersion(name, version string) (*ConfigVersion, error)

	// Create creates a new version of a configuration
	Create(config *ConfigVersion) error

	// Update updates an existing configuration version
	Update(config *ConfigVersion) error

	// List returns a list of configuration versions matching the filter
	List(filter ConfigVersionFilter) ([]*ConfigVersion, error)

	// Deactivate deactivates a specific version of a configuration
	Deactivate(name, version string) error
}
