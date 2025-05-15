package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aimustaev/service-workflow/internal/manager_workflow"
)

type ListNamesHandler struct {
	repo manager_workflow.ConfigVersionRepository
}

func NewListNamesHandler(repo manager_workflow.ConfigVersionRepository) *ListNamesHandler {
	return &ListNamesHandler{
		repo: repo,
	}
}

func (h *ListNamesHandler) Handle(w http.ResponseWriter, r *http.Request) {
	names, err := h.repo.ListNames()
	if err != nil {
		log.Printf("Error getting config names: %v", err)
		http.Error(w, "Failed to get config names", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(names); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
