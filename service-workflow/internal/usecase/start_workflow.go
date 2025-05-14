package usecase

import (
	"context"
	"log"

	"go.temporal.io/sdk/client"

	"github.com/aimustaev/service-workflow/internal/model"
	"github.com/aimustaev/service-workflow/internal/workflow"
)

type StartWorkflowInput struct {
	Message model.Message
}

type StartWorkflowOutput struct {
	WorkflowID string
}

type StartWorkflowUseCase struct {
	temporalClient client.Client
}

func NewStartWorkflowUseCase(temporalClient client.Client) *StartWorkflowUseCase {
	return &StartWorkflowUseCase{
		temporalClient: temporalClient,
	}
}

func (uc *StartWorkflowUseCase) Execute(ctx context.Context, message model.Message) (*StartWorkflowOutput, error) {
	workflowInput := workflow.SimpleWorkflowInput{
		Message: message,
	}

	workflowOptions := client.StartWorkflowOptions{
		ID:        "workflow-ticket-" + message.From,
		TaskQueue: "workflow-ticket",
	}

	log.Printf("Starting workflow execution with input: %s", workflowInput.Message)
	//workflowRun, err := uc.temporalClient.ExecuteWorkflow(
	//	ctx,
	//	workflowOptions,
	//	"SimpleWorkflow",
	//	workflowInput,
	//)

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

	return &StartWorkflowOutput{
		WorkflowID: workflowRun.GetID(),
	}, nil
}
