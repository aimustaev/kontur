package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aimustaev/service-workflow/internal/manager_workflow"
	"github.com/gorilla/mux"
)

type GetLatestConfigHandler struct {
	repo manager_workflow.ConfigVersionRepository
}

func NewGetLatestConfigHandler(repo manager_workflow.ConfigVersionRepository) *GetLatestConfigHandler {
	return &GetLatestConfigHandler{
		repo: repo,
	}
}

func (h *GetLatestConfigHandler) Handle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	if name == "" {
		http.Error(w, "Name parameter is required", http.StatusBadRequest)
		return
	}

	config, err := h.repo.GetLatestActive(name)
	if err != nil {
		log.Printf("Error getting latest active config: %v", err)
		http.Error(w, "Failed to get latest active config", http.StatusInternalServerError)
		return
	}
	if config == nil {
		http.Error(w, "Config not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config)
}
