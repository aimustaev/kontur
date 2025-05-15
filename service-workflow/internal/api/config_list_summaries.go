package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aimustaev/service-workflow/internal/manager_workflow"
)

type ListSummariesHandler struct {
	repo manager_workflow.ConfigVersionRepository
}

func NewListSummariesHandler(repo manager_workflow.ConfigVersionRepository) *ListSummariesHandler {
	return &ListSummariesHandler{
		repo: repo,
	}
}

func (h *ListSummariesHandler) Handle(w http.ResponseWriter, r *http.Request) {
	summaries, err := h.repo.ListSummaries()
	if err != nil {
		log.Printf("Error getting config summaries: %v", err)
		http.Error(w, "Failed to get config summaries", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(summaries); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
