package workflow

import (
	"fmt"
)

// WorkflowInput is the interface for all workflow inputs
type WorkflowInput interface {
	Validate() error
}

// WorkflowResult is the interface for all workflow results
type WorkflowResult interface {
	IsSuccess() bool
	GetError() error
}

// BaseWorkflowResult implements WorkflowResult interface
type BaseWorkflowResult struct {
	Success bool
	Err     error
}

// IsSuccess returns true if the result is successful
func (r BaseWorkflowResult) IsSuccess() bool {
	return r.Success
}

// GetError returns the error
func (r BaseWorkflowResult) GetError() error {
	return r.Err
}

// Error returns the error message
func (r BaseWorkflowResult) Error() string {
	if r.Err != nil {
		return fmt.Sprintf("workflow error: %v", r.Err)
	}
	return ""
}

// NewSuccessResult creates a new successful result
func NewSuccessResult() BaseWorkflowResult {
	return BaseWorkflowResult{
		Success: true,
	}
}

// NewErrorResult creates a new error result
func NewErrorResult(err error) BaseWorkflowResult {
	return BaseWorkflowResult{
		Success: false,
		Err:     err,
	}
}
