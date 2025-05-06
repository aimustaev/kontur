package workflow

import (
	"go.temporal.io/sdk/worker"
)

// RegisterWorkflows registers all workflows with the worker
func RegisterWorkflows(w worker.Worker) {
	// Register workflows
	w.RegisterWorkflow(SimpleWorkflow)
}
