package workflow

import (
	"time"

	"go.temporal.io/sdk/workflow"

	"github.com/aimustaev/service-workflow/internal/activity"
	"github.com/aimustaev/service-workflow/internal/generated/proto"
	"github.com/aimustaev/service-workflow/internal/model"
)

// SimpleWorkflowInput представляет входные данные для SimpleWorkflow
type SimpleWorkflowInput struct {
	Message model.Message
}

type TicketState struct {
	UserID   string
	Status   string
	Messages []model.Message
}

// SimpleWorkflow реализует функцию workflow
func SimpleWorkflow(ctx workflow.Context, input SimpleWorkflowInput) (BaseWorkflowResult, error) {
	message := input.Message

	state := TicketState{
		UserID:   "",
		Status:   "open",
		Messages: []model.Message{},
	}

	logger := workflow.GetLogger(ctx)

	logger.Info("Запуск workflow с входными данными", "message", message)

	// Настраиваем таймауты для активностей
	options := workflow.ActivityOptions{
		StartToCloseTimeout:    time.Minute * 5,  // 5 минут на выполнение активности
		ScheduleToCloseTimeout: time.Minute * 10, // 10 минут от планирования до завершения
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	// Выполняем обработку сообщения
	var processErr error
	err := workflow.ExecuteActivity(ctx, activity.ProcessMessageActivity, message.Body).Get(ctx, &processErr)
	if err != nil {
		return NewErrorResult(err), nil
	}

	state.Messages = append(state.Messages, message)

	// Выполняем обработку сообщения
	var ticket *proto.TicketResponse
	err = workflow.ExecuteActivity(ctx, activity.CreateTicketActivity, &proto.CreateTicketRequest{
		VerticalId: "",
		SkillId:    "",
		UserId:     message.From,
		Assign:     "",
	}).Get(ctx, &ticket)
	if err != nil {
		return NewErrorResult(err), nil
	}

	state.Status = "open"

	// Выполняем классификацию тикета
	var ticketByClassifier *model.Ticket
	err = workflow.ExecuteActivity(ctx, activity.ClassifierAcitivity, ticket).Get(ctx, &ticketByClassifier)
	if err != nil {
		return NewErrorResult(err), nil
	}

	state.Status = "classifier"

	// Выполняем ожидание
	var waitErr error
	err = workflow.ExecuteActivity(ctx, activity.WaitActivity, 10).Get(ctx, &waitErr)
	if err != nil {
		return NewErrorResult(err), nil
	}

	state.Status = "resolve"

	logger.Info("Workflow завершен")
	return NewSuccessResult(), nil
}
