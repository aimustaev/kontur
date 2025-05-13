package workflow

import (
	"time"

	"go.temporal.io/sdk/workflow"

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
func (w *Workflow) SimpleWorkflow(ctx workflow.Context, input SimpleWorkflowInput) (BaseWorkflowResult, error) {
	message := input.Message

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
	err := workflow.ExecuteActivity(ctx, w.activity.ProcessMessageActivity, message.Body).Get(ctx, &processErr)
	if err != nil {
		return NewErrorResult(err), nil
	}

	// Создаем тикет
	var ticket *proto.TicketResponse
	err = workflow.ExecuteActivity(ctx, w.activity.CreateTicketActivity, &proto.CreateTicketRequest{
		User:        message.From,
		Agent:       message.To,
		ProblemId:   0, // TODO: Добавить определение ProblemId
		VerticalId:  0, // TODO: Добавить определение VerticalId
		SkillId:     0, // TODO: Добавить определение SkillId
		UserGroupId: 0, // TODO: Добавить определение UserGroupId
		Channel:     message.Channel,
		Status:      "open",
	}).Get(ctx, &ticket)
	if err != nil {
		return NewErrorResult(err), nil
	}

	// Выполняем классификацию тикета
	var classifiedTicket *proto.TicketResponse
	err = workflow.ExecuteActivity(ctx, w.activity.ClassifierAcitivity, ticket).Get(ctx, &classifiedTicket)
	if err != nil {
		return NewErrorResult(err), nil
	}

	// Выполняем ожидание
	var waitErr error
	err = workflow.ExecuteActivity(ctx, w.activity.WaitActivity, 10).Get(ctx, &waitErr)
	if err != nil {
		return NewErrorResult(err), nil
	}

	logger.Info("Workflow завершен")
	return NewSuccessResult(classifiedTicket), nil
}
