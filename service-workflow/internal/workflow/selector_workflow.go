package workflow

import (
	"time"

	"go.temporal.io/sdk/workflow"

	"github.com/aimustaev/service-workflow/internal/activity"
	"github.com/aimustaev/service-workflow/internal/model"
)

// SimpleWorkflowInput представляет входные данные для SimpleWorkflow
type SelectorWorkflowInput struct {
	Message model.Message
}

func SelectorWorkflow(ctx workflow.Context, input SimpleWorkflowInput) (BaseWorkflowResult, error) {
	logger := workflow.GetLogger(ctx)

	// Настраиваем таймауты для активностей
	options := workflow.ActivityOptions{
		StartToCloseTimeout:    time.Minute * 5,  // 5 минут на выполнение активности
		ScheduleToCloseTimeout: time.Minute * 10, // 10 минут от планирования до завершения
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	logger.Info("Запуск workflow с входными данными", "message", input.Message)

	//workflow.Go(ctx, func(ctx workflow.Context) {
	//	ch := workflow.GetSignalChannel(ctx, "NewMessage")
	//	for {
	//		var msg model.Message
	//		ch.Receive(ctx, &msg) // Блокируется до получения сигнала
	//		// 3. Сохраняем сообщение в БД
	//		logger.Info("Запуск workflow с входными данными", "message", input.Message.Body)
	//	}
	//})

	err := workflow.ExecuteActivity(ctx, activity.WaitActivity, 30).Get(ctx, nil)
	if err != nil {
		return NewErrorResult(err), nil
	}

	return NewErrorResult(nil), nil
}
