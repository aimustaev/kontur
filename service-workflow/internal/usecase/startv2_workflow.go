package usecase

import (
	"context"
	"log"

	"go.temporal.io/sdk/client"

	"github.com/aimustaev/service-workflow/internal/model"
	"github.com/aimustaev/service-workflow/internal/workflow"
)

type StartV2WorkflowInput struct {
	Message model.Message
}

type StartV2WorkflowOutput struct {
	WorkflowID string
}

type StartV2WorkflowUseCase struct {
	temporalClient client.Client
}

func NewStartV2WorkflowUseCase(temporalClient client.Client) *StartV2WorkflowUseCase {
	return &StartV2WorkflowUseCase{
		temporalClient: temporalClient,
	}
}

func (uc *StartV2WorkflowUseCase) Execute(ctx context.Context, message model.Message) (*StartV2WorkflowOutput, error) {
	workflowInput := workflow.SelectorWorkflowInput{
		Message: message,
	}

	workflowOptions := client.StartWorkflowOptions{
		ID:        "workflow-ticket-" + message.From,
		TaskQueue: "workflow-ticket",
	}

	workflowRun, err := uc.temporalClient.SignalWithStartWorkflow(
		context.Background(),
		workflowOptions.ID, // ID Workflow (на основе userID)
		"NewMessage",       // Имя сигнала
		message,            // Данные сигнала
		workflowOptions,
		"DynamicTicketWorkflow", // Функция Workflow (если его нет)
		workflowInput,           // Аргументы Workflow
	)

	if err != nil {
		log.Printf("Error starting workflow: %v", err)
		return nil, err
	}

	return &StartV2WorkflowOutput{
		WorkflowID: workflowRun.GetID(),
	}, nil
}
