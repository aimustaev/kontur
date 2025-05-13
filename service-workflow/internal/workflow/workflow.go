package workflow

import (
	"go.temporal.io/sdk/worker"

	act "github.com/aimustaev/service-workflow/internal/activity"
	"github.com/aimustaev/service-workflow/internal/config"
)

// RegisterWorkflows registers all workflows with the worker
func RegisterWorkflows(w worker.Worker, cfg *config.Config) {
	// Set global config
	act.SetConfig(cfg)

	// Register workflows
	w.RegisterWorkflow(SimpleWorkflow)
	w.RegisterWorkflow(SelectorWorkflow)
	// Register activities
	w.RegisterActivity(act.CreateTicketActivity)
	w.RegisterActivity(act.ProcessMessageActivity)
	w.RegisterActivity(act.WaitActivity)
	w.RegisterActivity(act.ClassifierAcitivity)
}
