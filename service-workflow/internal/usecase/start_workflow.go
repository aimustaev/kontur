package usecase

import (
	"context"
	"log"
	"time"

	"go.temporal.io/sdk/client"

	"github.com/aimustaev/service-workflow/internal/workflow"
)

// StartWorkflowInput represents the input for starting a workflow
type StartWorkflowInput struct {
	Message string
}

// StartWorkflowOutput represents the output of starting a workflow
type StartWorkflowOutput struct {
	WorkflowID string
}

// StartWorkflowUseCase handles the business logic for starting workflows
type StartWorkflowUseCase struct {
	temporalClient client.Client
}

// NewStartWorkflowUseCase creates a new StartWorkflowUseCase
func NewStartWorkflowUseCase(temporalClient client.Client) *StartWorkflowUseCase {
	return &StartWorkflowUseCase{
		temporalClient: temporalClient,
	}
}

// Execute starts a new workflow
func (uc *StartWorkflowUseCase) Execute(ctx context.Context, input StartWorkflowInput) (*StartWorkflowOutput, error) {
	// Create workflow input
	workflowInput := workflow.SimpleWorkflowInput{
		Message: input.Message,
	}

	// Start workflow execution
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
