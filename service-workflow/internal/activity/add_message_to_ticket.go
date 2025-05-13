package activity

import (
	"context"
	"fmt"

	"go.temporal.io/sdk/activity"

	"github.com/aimustaev/service-workflow/internal/generated/proto"
	"github.com/aimustaev/service-workflow/internal/model"
)

// SaveMassageActivity сохраняет сообщение в тикет
func (a *Activity) AddMassageToTicketActivity(ctx context.Context, message *model.Message, ticketID string) (*proto.TicketResponse, error) {
	logger := activity.GetLogger(ctx)

	request := &proto.AddMessageToTicketRequest{
		TicketId:    ticketID,
		FromAddress: message.From,
		ToAddress:   message.To,
		Subject:     message.Subject,
		Body:        message.Body,
		Channel:     message.Channel,
	}
	// Создаем тикет
	_, err := a.ticketClient.AddMessageToTicket(ctx, request)
	if err != nil {
		logger.Error("Ошибка при добавлении сообщения в тикет", "error", err)
		return nil, fmt.Errorf("failed to add message to ticket: %w", err)
	}

	logger.Info("Сообщение успешно добавлено в тикет", "ticket_id", ticketID)
	return nil, nil
}
