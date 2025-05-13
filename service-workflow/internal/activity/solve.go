package activity

import (
	"context"
	"fmt"

	"go.temporal.io/sdk/activity"

	"github.com/aimustaev/service-workflow/internal/generated/proto"
)

// SolveTicketAcitivity решает тикет
func (a *Activity) SolveTicketAcitivity(ctx context.Context, request *proto.TicketResponse) (*proto.TicketResponse, error) {
	logger := activity.GetLogger(ctx)

	// Создаем тикет
	response, err := a.ticketClient.UpdateTicket(ctx, &proto.UpdateTicketRequest{
		Id:          request.Id,
		VerticalId:  request.VerticalId,
		ProblemId:   request.ProblemId,
		SkillId:     request.SkillId,
		UserGroupId: request.UserGroupId,
		User:        request.User,
		Agent:       request.Agent,
		Status:      "resolved",
		Channel:     request.Channel,
	})
	if err != nil {
		logger.Error("Ошибка при решении тикета", "error", err)
		return nil, fmt.Errorf("failed to solve ticket: %w", err)
	}

	logger.Info("Решение тикета", "ticket", response)

	return response, nil
}
