package activity

import (
	"context"
	"fmt"
	"github.com/aimustaev/service-workflow/internal/ticket"
	"math/rand"
	"strconv"

	"go.temporal.io/sdk/activity"

	"github.com/aimustaev/service-workflow/internal/generated/proto"
)

// ClassifierAcitivity классифицирует
func ClassifierAcitivity(ctx context.Context, request *proto.TicketResponse) (*proto.TicketResponse, error) {
	logger := activity.GetLogger(ctx)
	cfg := GetConfig()
	if cfg == nil {
		return nil, fmt.Errorf("config not found")
	}

	client, err := ticket.NewClient(cfg)
	if err != nil {
		logger.Error("Ошибка при создании клиента тикетов", "error", err)
		return nil, fmt.Errorf("failed to create ticket client: %w", err)
	}
	defer client.Close()

	// Создаем тикет
	response, err := client.UpdateTicket(ctx, &proto.UpdateTicketRequest{
		Id:         request.Id,
		VerticalId: strconv.Itoa(rand.Intn(10)),
		SkillId:    strconv.Itoa(rand.Intn(10)),
		UserId:     request.UserId,
		Assign:     request.Assign,
	})
	if err != nil {
		logger.Error("Ошибка при обновлении тикета", "error", err)
		return nil, fmt.Errorf("failed to update ticket: %w", err)
	}

	logger.Info("Классификация тикета", "ticket", response)

	return response, nil
}
