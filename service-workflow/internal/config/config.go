package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	HTTP     HTTPConfig
	Temporal TemporalConfig
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

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
