package usecase

import (
	"context"
	"log"

	"github.com/aimustaev/service-workflow/internal/model"
)

type Message struct {
	ID      string   `json:"id"`
	From    string   `json:"from"`
	To      string   `json:"to"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
	Tags    []string `json:"tags"`
	Channel string   `json:"channel"`
}

type MessageHandler struct {
	startWorkflowUC   StartWorkflowUseCase
	startV2WorkflowUC StartV2WorkflowUseCase
}

func NewMessageHandler(startWorkflowUC *StartWorkflowUseCase, startV2WorkflowUC *StartV2WorkflowUseCase) *MessageHandler {
	return &MessageHandler{
		startWorkflowUC:   *startWorkflowUC,
		startV2WorkflowUC: *startV2WorkflowUC,
	}
}

func (h *MessageHandler) HandleMessage(ctx context.Context, msg Message) error {
	log.Printf("Processing message: ID=%s, From=%s, To=%s, Channel=%s, %s",
		msg.ID, msg.From, msg.To, msg.Channel, msg.Body)

	h.startV2WorkflowUC.Execute(ctx, model.Message{
		ID:      msg.ID,
		From:    msg.From,
		To:      msg.To,
		Subject: msg.Subject,
		Body:    msg.Body,
		Tags:    msg.Tags,
		Channel: msg.Channel,
	})

	return nil
}
