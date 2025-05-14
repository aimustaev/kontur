package workflow

import (
	"fmt"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/workflow"

	act "github.com/aimustaev/service-workflow/internal/activity"
	"github.com/aimustaev/service-workflow/internal/config"
	"github.com/aimustaev/service-workflow/internal/engine"
)

type DynamicWorkflow struct {
	activity      *act.Activity
	engine        *engine.WorkflowEngine
	configManager *config.ConfigManager
}

func NewDynamicWorkflow(activity *act.Activity, temporalClient client.Client, configManager *config.ConfigManager) *DynamicWorkflow {
	// Создаем движок с маппингом activity
	activitiesMap := map[string]interface{}{
		"CreateTicketActivity":       activity.CreateTicketActivity,
		"ProcessMessageActivity":     activity.ProcessMessageActivity,
		"WaitActivity":               activity.WaitActivity,
		"ClassifierAcitivity":        activity.ClassifierAcitivity,
		"GetOrCreateTicketActivity":  activity.GetOrCreateTicketActivity,
		"GetTicketByUserActivity":    activity.GetTicketByUserActivity,
		"AddMassageToTicketActivity": activity.AddMassageToTicketActivity,
		"SolveTicketAcitivity":       activity.SolveTicketAcitivity,
	}

	return &DynamicWorkflow{
		activity:      activity,
		engine:        engine.NewEngine(temporalClient, activitiesMap),
		configManager: configManager,
	}
}

// DynamicWorkflow - универсальный обработчик для всех workflow
func (w *DynamicWorkflow) Execute(ctx workflow.Context, input map[string]interface{}) (interface{}, error) {
	// Получаем имя workflow из контекста
	workflowName := workflow.GetInfo(ctx).WorkflowType.Name

	// Получаем актуальное определение workflow из менеджера конфигураций
	def, err := w.configManager.GetWorkflowDefinition(workflowName)
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow definition: %w", err)
	}

	return w.engine.ExecuteWorkflow(ctx, def, input)
}
