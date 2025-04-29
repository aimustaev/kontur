package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// Config represents the application configuration
type Config struct {
	MailhogHost  string `env:"MAILHOG_HOST" default:"mailhog"`
	MailhogPort  int    `env:"MAILHOG_PORT" default:"8025"`
	LogLevel     string `env:"LOG_LEVEL" default:"info"`
	PostgresHost string `env:"POSTGRES_HOST" default:"postgres"`
	PostgresPort int    `env:"POSTGRES_PORT" default:"5432"`
	PostgresUser string `env:"POSTGRES_USER" default:"postgres"`
	PostgresPass string `env:"POSTGRES_PASSWORD" default:"postgres"`
	PostgresDB   string `env:"POSTGRES_DB" default:"service_adapter"`
}

// Load loads the configuration from environment variables
func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return &cfg, nil
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.PostgresUser,
		c.PostgresPass,
		c.PostgresHost,
		c.PostgresPort,
		c.PostgresDB,
	)
}
