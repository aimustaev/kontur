package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aimustaev/service-workflow/internal/manager_workflow"
	"github.com/gorilla/mux"
)

type UpdateConfigHandler struct {
	repo ConfigVersionRepository
}

func NewUpdateConfigHandler(repo ConfigVersionRepository) *UpdateConfigHandler {
	return &UpdateConfigHandler{
		repo: repo,
	}
}

func (h *UpdateConfigHandler) Handle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	version := vars["version"]
	if name == "" || version == "" {
		http.Error(w, "Name and version parameters are required", http.StatusBadRequest)
		return
	}

	var config manager_workflow.ConfigVersion
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	config.Name = name
	config.Version = version

	if err := h.repo.Update(&config); err != nil {
		log.Printf("Error updating config: %v", err)
		http.Error(w, "Failed to update config", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config)
}
