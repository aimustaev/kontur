package temporal

import (
	"log"

	"go.temporal.io/sdk/client"
)

// Config holds the configuration for Temporal client
type Config struct {
	HostPort  string
	Namespace string
}

// DefaultConfig returns default configuration
func DefaultConfig() Config {
	return Config{
		HostPort:  "localhost:7233",
		Namespace: "default",
	}
}

// Client wraps Temporal client with additional functionality
type Client struct {
	client client.Client
	config Config
}

// NewClient creates a new Temporal client
func NewClient(config Config) (*Client, error) {
	log.Println("Initializing Temporal client...")
	c, err := client.NewClient(client.Options{
		HostPort:  config.HostPort,
		Namespace: config.Namespace,
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		client: c,
		config: config,
	}, nil
}

// Close closes the Temporal client
func (c *Client) Close() {
	c.client.Close()
}

// GetClient returns the underlying Temporal client
func (c *Client) GetClient() client.Client {
	return c.client
}

// GetConfig returns the client configuration
func (c *Client) GetConfig() Config {
	return c.config
}
