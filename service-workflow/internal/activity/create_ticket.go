package activity

import (
	"context"
	"fmt"

	"go.temporal.io/sdk/activity"

	"github.com/aimustaev/service-workflow/internal/generated/proto"
)

// CreateTicketActivity создает новый тикет через gRPC сервис
func (a *Activity) CreateTicketActivity(ctx context.Context, request *proto.CreateTicketRequest) (*proto.TicketResponse, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Создание нового тикета", "user", request.User)

	// Создаем тикет
	response, err := a.ticketClient.CreateTicket(ctx, request)
	if err != nil {
		logger.Error("Ошибка при создании тикета", "error", err)
		return nil, fmt.Errorf("failed to create ticket: %w", err)
	}

	logger.Info("Тикет успешно создан", "ticket_id", response.Id)
	return response, nil
}
