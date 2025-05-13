package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aimustaev/service-workflow/internal/model"
	"github.com/aimustaev/service-workflow/internal/usecase"
)

type StartV2WorkflowRequest struct {
	ID      string   `json:"id"`
	From    string   `json:"from"`
	To      string   `json:"to"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
	Tags    []string `json:"tags"`
	Channel string   `json:"channel"` // email, telegram
}

type StartV2WorkflowResponse struct {
	WorkflowID string `json:"workflow_id"`
}

type StartV2WorkflowHandler struct {
	startWorkflowUseCase *usecase.StartV2WorkflowUseCase
}

func NewStartV2WorkflowHandler(startWorkflowUseCase *usecase.StartV2WorkflowUseCase) *StartV2WorkflowHandler {
	return &StartV2WorkflowHandler{
		startWorkflowUseCase: startWorkflowUseCase,
	}
}

func (h *StartV2WorkflowHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req StartV2WorkflowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	output, err := h.startWorkflowUseCase.Execute(r.Context(), model.Message{
		ID:      req.ID,
		From:    req.From,
		To:      req.To,
		Subject: req.Subject,
		Body:    req.Body,
		Tags:    req.Tags,
		Channel: req.Channel},
	)

	if err != nil {
		log.Printf("Error executing usecase: %v", err)
		http.Error(w, "Failed to start workflow", http.StatusInternalServerError)
		return
	}

	// Return workflow ID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(StartWorkflowResponse{
		WorkflowID: output.WorkflowID,
	})
}
