package workflow

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

// SimpleWorkflowInput represents the input for SimpleWorkflow
type SimpleWorkflowInput struct {
	Message string
}

func (i SimpleWorkflowInput) Validate() error {
	if i.Message == "" {
		return ErrEmptyMessage
	}
	return nil
}

// SimpleWorkflow implements the workflow function
func SimpleWorkflow(ctx workflow.Context, input SimpleWorkflowInput) (BaseWorkflowResult, error) {
	logger := workflow.GetLogger(ctx)

	// Validate input
	if err := input.Validate(); err != nil {
		return NewErrorResult(err), nil
	}

	logger.Info("Starting workflow with input", "message", input.Message)

	// Use workflow.Sleep instead of time.Sleep
	logger.Info("Sleeping for 5 seconds...")
	workflow.Sleep(ctx, 5*time.Second)

	logger.Info("Workflow completed")
	return NewSuccessResult(), nil
}
