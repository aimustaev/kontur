package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	HTTP     HTTPConfig
	Temporal TemporalConfig
	Kafka    KafkaConfig
	Ticket   TicketConfig
	Postgres PostgresConfig
}

// HTTPConfig holds HTTP server configuration
type HTTPConfig struct {
	Host string
	Port string
}

// TemporalConfig holds Temporal client configuration
type TemporalConfig struct {
	Host      string
	Port      string
	Namespace string
}

// KafkaConfig holds Kafka configuration
type KafkaConfig struct {
	Brokers []string
	GroupID string
	Topic   string
}

// TicketConfig holds ticket service configuration
type TicketConfig struct {
	Host string
	Port string
}

// PostgresConfig holds PostgreSQL configuration
type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

// Load loads configuration from environment variables
func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		HTTP: HTTPConfig{
			Host: getEnv("HTTP_HOST", "0.0.0.0"),
			Port: getEnv("HTTP_PORT", "3002"),
		},
		Temporal: TemporalConfig{
			Host:      getEnv("TEMPORAL_HOST", "localhost"),
			Port:      getEnv("TEMPORAL_PORT", "7233"),
			Namespace: getEnv("TEMPORAL_NAMESPACE", "default"),
		},
		Kafka: KafkaConfig{
			Brokers: []string{getEnv("KAFKA_BROKERS", "kafka:9092")},
			GroupID: getEnv("KAFKA_GROUP_ID", "service-workflow-group"),
			Topic:   getEnv("KAFKA_TOPIC", "workflow-topic"),
		},
		Ticket: TicketConfig{
			Host: getEnv("TICKET_SERVICE_HOST", "localhost"),
			Port: getEnv("TICKET_SERVICE_PORT", "50051"),
		},
		Postgres: PostgresConfig{
			Host:     getEnv("POSTGRES_HOST", "localhost"),
			Port:     getEnv("POSTGRES_PORT", "5432"),
			User:     getEnv("POSTGRES_USER", "postgres"),
			Password: getEnv("POSTGRES_PASSWORD", "postgres"),
			Database: getEnv("POSTGRES_DB", "service_tickets"),
			SSLMode:  getEnv("POSTGRES_SSLMODE", "disable"),
		},
	}
}

// GetHTTPAddr returns the full HTTP address
func (c *Config) GetHTTPAddr() string {
	return c.HTTP.Host + ":" + c.HTTP.Port
}

// GetTemporalAddr returns the full Temporal address
func (c *Config) GetTemporalAddr() string {
	return c.Temporal.Host + ":" + c.Temporal.Port
}

// GetTicketServiceAddr returns the full ticket service address
func (c *Config) GetTicketServiceAddr() string {
	return c.Ticket.Host + ":" + c.Ticket.Port
}

// GetPostgresDSN returns the PostgreSQL connection string
func (c *Config) GetPostgresDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Postgres.Host,
		c.Postgres.Port,
		c.Postgres.User,
		c.Postgres.Password,
		c.Postgres.Database,
		c.Postgres.SSLMode,
	)
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
