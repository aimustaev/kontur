package mailhog

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/aimustaev/service-adapter/internal/adapter"
)

// MailhogAdapter implements the Adapter interface for Mailhog
type MailhogAdapter struct {
	host      string
	port      int
	client    *http.Client
	logger    *logrus.Logger
	processed map[string]bool
	db        adapter.Database
}

// NewMailhogAdapter creates a new Mailhog adapter
func NewMailhogAdapter(host string, port int, logger *logrus.Logger, db adapter.Database) *MailhogAdapter {
	return &MailhogAdapter{
		host:      host,
		port:      port,
		client:    &http.Client{Timeout: 10 * time.Second},
		logger:    logger,
		processed: make(map[string]bool),
		db:        db,
	}
}

// Connect establishes a connection to Mailhog
func (a *MailhogAdapter) Connect(ctx context.Context) error {
	// Mailhog is HTTP-based, so we just need to check if it's accessible
	url := fmt.Sprintf("http://%s:%d/api/v2/messages", a.host, a.port)
	a.logger.Infof("Connecting to Mailhog at %s", url)

	resp, err := a.client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to connect to Mailhog: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code from Mailhog: %d", resp.StatusCode)
	}

	return nil
}

// Disconnect closes the connection to Mailhog
func (a *MailhogAdapter) Disconnect(ctx context.Context) error {
	// No connection to close for HTTP-based service
	return nil
}

// GetMessages retrieves messages from Mailhog
func (a *MailhogAdapter) GetMessages(ctx context.Context) ([]adapter.Message, error) {
	url := fmt.Sprintf("http://%s:%d/api/v2/messages", a.host, a.port)
	a.logger.Infof("Getting messages from %s", url)

	resp, err := a.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var data struct {
		Items []struct {
			ID      string `json:"ID"`
			Content struct {
				Headers struct {
					From    []string `json:"From"`
					To      []string `json:"To"`
					Subject []string `json:"Subject"`
				} `json:"Headers"`
				Body string `json:"Body"`
			} `json:"Content"`
			Tags []string `json:"Tags"`
		} `json:"items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	var messages []adapter.Message
	for _, item := range data.Items {
		// Skip already processed messages
		if a.processed[item.ID] {
			continue
		}

		message := adapter.Message{
			ID:      item.ID,
			From:    getFirstHeader(item.Content.Headers.From),
			To:      getFirstHeader(item.Content.Headers.To),
			Subject: getFirstHeader(item.Content.Headers.Subject),
			Body:    item.Content.Body,
			Tags:    []string{"processed"},
		}
		messages = append(messages, message)
		// Save email to database
		if err := a.db.SaveEmail(ctx, message); err != nil {
			a.logger.Errorf("Failed to save email %s to database: %v", message.ID, err)
		}
	}

	return messages, nil
}

// MarkAsProcessed marks a message as processed
func (a *MailhogAdapter) MarkAsProcessed(ctx context.Context, messageID string) error {
	a.processed[messageID] = true
	a.logger.Infof("Marked message %s as processed", messageID)
	return nil
}

// getFirstHeader returns the first value from a header array or empty string
func getFirstHeader(headers []string) string {
	if len(headers) > 0 {
		return headers[0]
	}
	return ""
}
