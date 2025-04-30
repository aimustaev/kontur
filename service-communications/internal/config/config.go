package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// Config represents the application configuration
type Config struct {
	// Mailhog configuration
	MailhogHost string `envconfig:"MAILHOG_HOST" default:"mailhog"`
	MailhogPort int    `envconfig:"MAILHOG_PORT" default:"8025"`

	// Telegram configuration
	TelegramBotToken string `envconfig:"TELEGRAM_BOT_TOKEN"`
	TelegramChatID   string `envconfig:"TELEGRAM_CHAT_ID"`

	// Database configuration
	PostgresHost string `envconfig:"POSTGRES_HOST" default:"postgres"`
	PostgresPort int    `envconfig:"POSTGRES_PORT" default:"5432"`
	PostgresUser string `envconfig:"POSTGRES_USER" default:"postgres"`
	PostgresPass string `envconfig:"POSTGRES_PASSWORD"`
	PostgresDB   string `envconfig:"POSTGRES_DB" default:"service_communications"`

	// Gateway configuration
	GatewayHost string `envconfig:"GATEWAY_HOST" default:"service-gateway"`
	GatewayPort int    `envconfig:"GATEWAY_PORT" default:"50051"`

	// Logging configuration
	LogLevel string `envconfig:"LOG_LEVEL" default:"info"`
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

// GetGatewayAddress returns the full gateway address
func (c *Config) GetGatewayAddress() string {
	return fmt.Sprintf("%s:%d", c.GatewayHost, c.GatewayPort)
}
