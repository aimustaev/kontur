package activity

import (
	"context"
	"fmt"

	"go.temporal.io/sdk/activity"

	"github.com/aimustaev/service-workflow/internal/generated/proto"
	"github.com/aimustaev/service-workflow/internal/ticket"
)

// CreateTicketActivity создает новый тикет через gRPC сервис
func CreateTicketActivity(ctx context.Context, request *proto.CreateTicketRequest) (*proto.TicketResponse, error) {
	cfg := GetConfig()
	if cfg == nil {
		return nil, fmt.Errorf("config not found")
	}

	logger := activity.GetLogger(ctx)
	logger.Info("Создание нового тикета",
		"vertical_id", request.VerticalId,
		"user_id", request.UserId,
		"assign", request.Assign,
		"skill_id", request.SkillId,
	)

	// Создаем клиент для работы с сервисом тикетов
	client, err := ticket.NewClient(cfg)
	if err != nil {
		logger.Error("Ошибка при создании клиента тикетов", "error", err)
		return nil, fmt.Errorf("failed to create ticket client: %w", err)
	}
	defer client.Close()

	// Создаем тикет
	response, err := client.CreateTicket(ctx, request)
	if err != nil {
		logger.Error("Ошибка при создании тикета", "error", err)
		return nil, fmt.Errorf("failed to create ticket: %w", err)
	}

	logger.Info("Тикет успешно создан", "ticket_id", response.Id)
	return response, nil
}
