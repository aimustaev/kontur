package workflow

import (
	"go.temporal.io/sdk/worker"

	act "github.com/aimustaev/service-workflow/internal/activity"
	"github.com/aimustaev/service-workflow/internal/generated/proto"
)

type Workflow struct {
	activity *act.Activity
}

func NewWorkflow(activity *act.Activity) *Workflow {
	return &Workflow{
		activity: activity,
	}
}

// RegisterWorkflows registers all workflows with the worker
func RegisterWorkflows(w worker.Worker, ticketClient proto.TicketServiceClient) {
	activity := act.NewActivity(ticketClient)
	workflow := NewWorkflow(activity)

	// Register workflows
	w.RegisterWorkflow(workflow.SimpleWorkflow)
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
