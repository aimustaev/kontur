package workflow

import (
	"go.temporal.io/sdk/workflow"

	"github.com/aimustaev/service-workflow/internal/activity"
)

// SimpleWorkflowInput представляет входные данные для SimpleWorkflow
type SimpleWorkflowInput struct {
	Message string
}

func (i SimpleWorkflowInput) Validate() error {
	if i.Message == "" {
		return ErrEmptyMessage
	}
	return nil
}

// SimpleWorkflow реализует функцию workflow
func SimpleWorkflow(ctx workflow.Context, input SimpleWorkflowInput) (BaseWorkflowResult, error) {
	logger := workflow.GetLogger(ctx)

	// Валидация входных данных
	if err := input.Validate(); err != nil {
		return NewErrorResult(err), nil
	}

	logger.Info("Запуск workflow с входными данными", "message", input.Message)

	// Выполняем обработку сообщения
	var processErr error
	err := workflow.ExecuteActivity(ctx, activity.ProcessMessageActivity, input.Message).Get(ctx, &processErr)
	if err != nil {
		return NewErrorResult(err), nil
	}

	// Выполняем ожидание
	var waitErr error
	err = workflow.ExecuteActivity(ctx, activity.WaitActivity, 5).Get(ctx, &waitErr)
	if err != nil {
		return NewErrorResult(err), nil
	}

	logger.Info("Workflow завершен")
	return NewSuccessResult(), nil
}
