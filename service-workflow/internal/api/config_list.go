package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aimustaev/service-workflow/internal/manager_workflow"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type ListConfigHandler struct {
	repo manager_workflow.ConfigVersionRepository
}

func NewListConfigHandler(repo manager_workflow.ConfigVersionRepository) *ListConfigHandler {
	return &ListConfigHandler{
		repo: repo,
	}
}

func (h *ListConfigHandler) Handle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	if idStr == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	filter := manager_workflow.ConfigVersionFilter{
		ID: &id,
	}

	configs, err := h.repo.List(filter)
	if err != nil {
		log.Printf("Error listing configs: %v", err)
		http.Error(w, "Failed to list configs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(configs)
}
