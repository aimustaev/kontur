package workflow

import (
	"fmt"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/workflow"

	act "github.com/aimustaev/service-workflow/internal/activity"
	"github.com/aimustaev/service-workflow/internal/engine"
)

type DynamicWorkflow struct {
	activity     *act.Activity
	engine       *engine.WorkflowEngine
	workflowDefs map[string]engine.WorkflowDefinition // Кеш определений
}

func NewDynamicWorkflow(activity *act.Activity, temporalClient client.Client) *DynamicWorkflow {
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
		activity:     activity,
		engine:       engine.NewEngine(temporalClient, activitiesMap),
		workflowDefs: make(map[string]engine.WorkflowDefinition),
	}
}

// DynamicWorkflow - универсальный обработчик для всех workflow
func (w *DynamicWorkflow) Execute(ctx workflow.Context, input map[string]interface{}) (interface{}, error) {
	// Получаем имя workflow из контекста
	workflowName := workflow.GetInfo(ctx).WorkflowType.Name

	// Загружаем определение (в реальности можно из БД/файла/кэша)
	def, ok := w.workflowDefs[workflowName]
	if !ok {
		return nil, fmt.Errorf("workflow definition not found: %s", workflowName)
	}

	return w.engine.ExecuteWorkflow(ctx, def, input)
}

// AddDefinition добавляет/обновляет определение workflow
func (w *DynamicWorkflow) AddDefinition(name string, def engine.WorkflowDefinition) {
	w.workflowDefs[name] = def
}
