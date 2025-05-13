package activity

import (
	"context"
	"fmt"

	"go.temporal.io/sdk/activity"

	"github.com/aimustaev/service-workflow/internal/generated/proto"
	"github.com/aimustaev/service-workflow/internal/model"
)

// GetOrCreateTicketActivity получает тикет по пользователю
func (a *Activity) GetOrCreateTicketActivity(ctx context.Context, message *model.Message) (*proto.TicketResponse, error) {
	logger := activity.GetLogger(ctx)

	ticket, err := a.GetTicketByUserActivity(ctx, message.From)
	if err != nil {
		logger.Error("Ошибка при получении тикета", "error", err)
	}

	if ticket != nil {
		logger.Info("Тикет успешно получен", "ticket_id", ticket.Id)
		return ticket, nil
	}

	ticket, err = a.CreateTicketActivity(ctx, &proto.CreateTicketRequest{
		User:    message.From,
		Channel: message.Channel,
		Status:  "new",
	})
	if err != nil {
		logger.Error("Ошибка при создании тикета", "error", err)
		return nil, fmt.Errorf("failed to create ticket: %w", err)
	}

	logger.Info("Тикет успешно создан", "ticket_id", ticket.Id)

	return ticket, nil
}
