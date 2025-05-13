package workflow

import (
	"time"

	"go.temporal.io/sdk/workflow"

	"github.com/aimustaev/service-workflow/internal/generated/proto"
	"github.com/aimustaev/service-workflow/internal/model"
)

// SimpleWorkflowInput представляет входные данные для SimpleWorkflow
type SelectorWorkflowInput struct {
	Message model.Message
}

type SelectorWorkflowState struct {
	TicketID string
}

func (w *Workflow) SelectorWorkflow(ctx workflow.Context, input SimpleWorkflowInput) (BaseWorkflowResult, error) {
	logger := workflow.GetLogger(ctx)

	// Настраиваем таймауты для активностей
	options := workflow.ActivityOptions{
		StartToCloseTimeout:    time.Minute * 5,  // 5 минут на выполнение активности
		ScheduleToCloseTimeout: time.Minute * 10, // 10 минут от планирования до завершения
	}
	ctx = workflow.WithActivityOptions(ctx, options)
	logger.Info("Запуск workflow с входными данными", "message", input.Message)

	// Шаг 1. получение или создание тикета
	var ticket *proto.TicketResponse
	err := workflow.ExecuteActivity(ctx, w.activity.GetOrCreateTicketActivity, &input.Message).Get(ctx, &ticket)
	if err != nil {
		return NewErrorResult(err), nil
	}

	// Шаг 2. сохранение сообщения в тикет
	err = workflow.ExecuteActivity(ctx, w.activity.AddMassageToTicketActivity, &input.Message, ticket.Id).Get(ctx, nil)
	if err != nil {
		return NewErrorResult(err), nil
	}

	workflow.Go(ctx, func(ctx workflow.Context) {
		msgCh := workflow.GetSignalChannel(ctx, "NewMessage")
		for {
			var msg model.Message
			msgCh.Receive(ctx, &msg) // Ждём новые сообщения
			_ = workflow.ExecuteActivity(ctx, w.activity.AddMassageToTicketActivity, msg, ticket.Id).Get(ctx, nil)
			// 3. Сохраняем сообщения в БД
			logger.Info("Новое сообщение", "message", input.Message.Body, "time", time.Now().Format(time.RFC3339))
		}
	})

	// 4. Классификация тикета
	err = workflow.ExecuteActivity(ctx, w.activity.ClassifierAcitivity, ticket).Get(ctx, &ticket)
	if err != nil {
		return NewErrorResult(err), nil
	}

	// 5. Ожидание
	// TODO: Добавить ожидание
	err = workflow.ExecuteActivity(ctx, w.activity.WaitActivity, 30).Get(ctx, nil)
	if err != nil {
		return NewErrorResult(err), nil
	}

	// 6. Решение тикета
	// TODO: Добавить решение тикета
	err = workflow.ExecuteActivity(ctx, w.activity.SolveTicketAcitivity, ticket).Get(ctx, &ticket)
	if err != nil {
		return NewErrorResult(err), nil
	}

	return NewSuccessResult(ticket), nil
}
