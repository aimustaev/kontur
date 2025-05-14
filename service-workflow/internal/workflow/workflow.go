package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	workflow2 "go.temporal.io/sdk/workflow"

	act "github.com/aimustaev/service-workflow/internal/activity"
	"github.com/aimustaev/service-workflow/internal/config"
	"github.com/aimustaev/service-workflow/internal/engine"
	"github.com/aimustaev/service-workflow/internal/generated/proto"
)

type Workflow struct {
	activity   *act.Activity
	configRepo config.ConfigVersionRepository
}

func NewWorkflow(activity *act.Activity, configRepo config.ConfigVersionRepository) *Workflow {
	return &Workflow{
		activity:   activity,
		configRepo: configRepo,
	}
}

// RegisterWorkflows registers all workflows with the worker
func RegisterWorkflows(w worker.Worker, ticketClient proto.TicketServiceClient, temporalClient client.Client, configRepo config.ConfigVersionRepository) {
	activity := act.NewActivity(ticketClient)
	workflow := NewWorkflow(activity, configRepo)

	dynamicWorkflow := NewDynamicWorkflow(activity, temporalClient)

	// Load workflow definition from database
	workflowDef, err := loadWorkflowDefinitionFromDB(context.Background(), configRepo, "SimpleWorkflow")
	if err != nil {
		panic(fmt.Errorf("failed to load workflow definition: %w", err))
	}

	dynamicWorkflow.AddDefinition("DynamicTicketWorkflow", workflowDef)

	// Register workflows
	w.RegisterWorkflowWithOptions(dynamicWorkflow.Execute, workflow2.RegisterOptions{Name: "DynamicTicketWorkflow"})
	w.RegisterWorkflow(workflow.SelectorWorkflow)

	// Register activities
	w.RegisterActivity(workflow.activity.CreateTicketActivity)
	w.RegisterActivity(workflow.activity.ProcessMessageActivity)
	w.RegisterActivity(workflow.activity.WaitActivity)
	w.RegisterActivity(workflow.activity.ClassifierAcitivity)
	w.RegisterActivity(workflow.activity.GetOrCreateTicketActivity)
	w.RegisterActivity(workflow.activity.GetTicketByUserActivity)
	w.RegisterActivity(workflow.activity.AddMassageToTicketActivity)
	w.RegisterActivity(workflow.activity.SolveTicketAcitivity)
}

// loadWorkflowDefinitionFromDB loads a workflow definition from the database
func loadWorkflowDefinitionFromDB(ctx context.Context, repo config.ConfigVersionRepository, name string) (engine.WorkflowDefinition, error) {
	config, err := repo.GetLatestActive(name)
	if err != nil {
		return engine.WorkflowDefinition{}, fmt.Errorf("failed to get latest active config: %w", err)
	}
	if config == nil {
		return engine.WorkflowDefinition{}, fmt.Errorf("no active configuration found for workflow: %s", name)
	}

	var def engine.WorkflowDefinition
	if err := json.Unmarshal(config.Content, &def); err != nil {
		return engine.WorkflowDefinition{}, fmt.Errorf("failed to unmarshal workflow definition: %w", err)
	}

	return def, nil
}

// loadWorkflowDefinition is kept for backward compatibility and testing
func loadWorkflowDefinition(filename string) engine.WorkflowDefinition {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	var def engine.WorkflowDefinition
	if err := json.Unmarshal(data, &def); err != nil {
		panic(err)
	}

	return def
}
