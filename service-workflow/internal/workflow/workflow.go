package workflow

import (
	"context"
	"github.com/aimustaev/service-workflow/internal/manager_workflow"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	workflow2 "go.temporal.io/sdk/workflow"

	act "github.com/aimustaev/service-workflow/internal/activity"
	"github.com/aimustaev/service-workflow/internal/generated/proto"
)

type Workflow struct {
	activity      *act.Activity
	configManager *manager_workflow.ConfigManager
}

func NewWorkflow(activity *act.Activity, configManager *manager_workflow.ConfigManager) *Workflow {
	return &Workflow{
		activity:      activity,
		configManager: configManager,
	}
}

// RegisterWorkflows registers all workflows with the worker
func RegisterWorkflows(w worker.Worker, ticketClient proto.TicketServiceClient, temporalClient client.Client, configRepo manager_workflow.ConfigVersionRepository) {
	activity := act.NewActivity(ticketClient)

	// Создаем менеджер конфигураций с интервалом обновления 1 минута
	configManager := manager_workflow.NewConfigManager(configRepo, time.Minute)
	configManager.Start(context.Background())

	workflow := NewWorkflow(activity, configManager)
	dynamicWorkflow := NewDynamicWorkflow(activity, temporalClient, configManager)

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
