package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aimustaev/service-workflow/internal/usecase"
)

// StartWorkflowRequest represents the request body for starting a workflow
type StartWorkflowRequest struct {
	Message string `json:"message"`
}

// StartWorkflowResponse represents the response for starting a workflow
type StartWorkflowResponse struct {
	WorkflowID string `json:"workflow_id"`
}

// StartWorkflowHandler handles the /start endpoint
type StartWorkflowHandler struct {
	startWorkflowUseCase *usecase.StartWorkflowUseCase
}

// NewStartWorkflowHandler creates a new StartWorkflowHandler
func NewStartWorkflowHandler(startWorkflowUseCase *usecase.StartWorkflowUseCase) *StartWorkflowHandler {
	return &StartWorkflowHandler{
		startWorkflowUseCase: startWorkflowUseCase,
	}
}

// Handle handles the /start endpoint
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

	// Execute usecase
	output, err := h.startWorkflowUseCase.Execute(r.Context(), usecase.StartWorkflowInput{
		Message: req.Message,
	})
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
