package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aimustaev/service-workflow/internal/manager_workflow"
	"github.com/gorilla/mux"
)

type ListConfigHandler struct {
	repo ConfigVersionRepository
}

func NewListConfigHandler(repo ConfigVersionRepository) *ListConfigHandler {
	return &ListConfigHandler{
		repo: repo,
	}
}

func (h *ListConfigHandler) Handle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	if name == "" {
		http.Error(w, "Name parameter is required", http.StatusBadRequest)
		return
	}

	filter := manager_workflow.ConfigVersionFilter{
		Name: &name,
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
