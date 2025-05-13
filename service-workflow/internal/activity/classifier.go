package activity

import (
	"context"
	crand "crypto/rand"
	"fmt"
	"math/big"

	"go.temporal.io/sdk/activity"

	"github.com/aimustaev/service-workflow/internal/generated/proto"
)

// ClassifierAcitivity классифицирует
func (a *Activity) ClassifierAcitivity(ctx context.Context, request *proto.TicketResponse) (*proto.TicketResponse, error) {
	logger := activity.GetLogger(ctx)

	// Генерируем случайные числа для VerticalId и SkillId
	verticalId, err := crand.Int(crand.Reader, big.NewInt(10))
	if err != nil {
		return nil, fmt.Errorf("failed to generate vertical id: %w", err)
	}
	skillId, err := crand.Int(crand.Reader, big.NewInt(10))
	if err != nil {
		return nil, fmt.Errorf("failed to generate skill id: %w", err)
	}

	problemId, err := crand.Int(crand.Reader, big.NewInt(10))
	if err != nil {
		return nil, fmt.Errorf("failed to generate problem id: %w", err)
	}
	userGroupId, err := crand.Int(crand.Reader, big.NewInt(10))
	if err != nil {
		return nil, fmt.Errorf("failed to generate user group id: %w", err)
	}

	// Создаем тикет
	response, err := a.ticketClient.UpdateTicket(ctx, &proto.UpdateTicketRequest{
		Id:          request.Id,
		VerticalId:  verticalId.Int64(),
		ProblemId:   problemId.Int64(),
		SkillId:     skillId.Int64(),
		UserGroupId: userGroupId.Int64(),
		User:        request.User,
		Agent:       request.Agent,
		Status:      request.Status,
		Channel:     request.Channel,
	})
	if err != nil {
		logger.Error("Ошибка при обновлении тикета", "error", err)
		return nil, fmt.Errorf("failed to update ticket: %w", err)
	}

	logger.Info("Классификация тикета", "ticket", response)

	return response, nil
}
