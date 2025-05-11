package usecase

import (
	"context"
	"log"
)

// Message represents the structure of a message
type Message struct {
	ID      string   `json:"id"`
	From    string   `json:"from"`
	To      string   `json:"to"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
	Tags    []string `json:"tags"`
	Channel string   `json:"channel"` // email, telegram
}

// MessageHandler represents a usecase for handling Kafka messages
type MessageHandler struct {
	startWorkflowUC StartWorkflowUseCase
	// Add any dependencies here if needed
	// For example: repository, service clients, etc.
}

// NewMessageHandler creates a new message handler usecase
func NewMessageHandler(startWorkflowUC *StartWorkflowUseCase) *MessageHandler {
	return &MessageHandler{
		startWorkflowUC: *startWorkflowUC,
	}
}

// HandleMessage processes a message
func (h *MessageHandler) HandleMessage(ctx context.Context, msg Message) error {
	// Log the received message
	log.Printf("Processing message: ID=%s, From=%s, To=%s, Channel=%s, %s",
		msg.ID, msg.From, msg.To, msg.Channel, msg.Body)

	// Start workflow
	h.startWorkflowUC.Execute(ctx, StartWorkflowInput{
		Message: msg.Body,
	})
	// Here you can add your business logic for processing the message
	// For example:
	// - Validate the message data
	// - Process based on channel type
	// - Store in database
	// - Call other services
	// - etc.

	return nil
}
