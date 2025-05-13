package activity

import (
	"context"
	"fmt"

	"go.temporal.io/sdk/activity"

	"github.com/aimustaev/service-workflow/internal/generated/proto"
)

// GetTicketByUserActivity получает тикет по пользователю
func (a *Activity) GetTicketByUserActivity(ctx context.Context, user string) (*proto.TicketResponse, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Получение тикета по пользователю", "user", user)

	// Создаем тикет
	response, err := a.ticketClient.GetActiveTicketsByUser(ctx, &proto.GetActiveTicketsByUserRequest{
		User: user,
	})

	if err != nil {
		logger.Error("Ошибка при получении тикета", "error", err)
		return nil, fmt.Errorf("failed to get ticket: %w", err)
	}

	if len(response.Tickets) == 0 {
		logger.Info("Тикет не найден", "user", user)
		return nil, fmt.Errorf("ticket not found")
	}

	ticket := response.Tickets[0]

	logger.Info("Тикет успешно получен", "ticket_id", ticket.Id)
	return ticket, nil
}
