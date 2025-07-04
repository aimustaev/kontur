package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aimustaev/service-workflow/internal/model"
	"github.com/aimustaev/service-workflow/internal/usecase"
)

type StartWorkflowRequest struct {
	ID      string   `json:"id"`
	From    string   `json:"from"`
	To      string   `json:"to"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
	Tags    []string `json:"tags"`
	Channel string   `json:"channel"` // email, telegram
}

type StartWorkflowResponse struct {
	WorkflowID string `json:"workflow_id"`
}

type StartWorkflowHandler struct {
	startWorkflowUseCase *usecase.StartWorkflowUseCase
}

func NewStartWorkflowHandler(startWorkflowUseCase *usecase.StartWorkflowUseCase) *StartWorkflowHandler {
	return &StartWorkflowHandler{
		startWorkflowUseCase: startWorkflowUseCase,
	}
}

func (h *StartWorkflowHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req StartWorkflowRequest
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
