package usecase

import (
	"context"
	"log"
	"time"

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
		ID:        "workflow-ticket-" + time.Now().Format("20060102150405"),
		TaskQueue: "workflow-ticket",
	}

	log.Printf("Starting workflow execution with input: %s", workflowInput.Message)
	workflowRun, err := uc.temporalClient.ExecuteWorkflow(
		ctx,
		workflowOptions,
		"SimpleWorkflow",
		workflowInput,
	)
	if err != nil {
		log.Printf("Error starting workflow: %v", err)
		return nil, err
	}

	return &StartWorkflowOutput{
		WorkflowID: workflowRun.GetID(),
	}, nil
}
