package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aimustaev/service-workflow/internal/manager_workflow"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type UpdateConfigHandler struct {
	repo manager_workflow.ConfigVersionRepository
}

func NewUpdateConfigHandler(repo manager_workflow.ConfigVersionRepository) *UpdateConfigHandler {
	return &UpdateConfigHandler{
		repo: repo,
	}
}

func (h *UpdateConfigHandler) Handle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	version := vars["version"]
	if idStr == "" || version == "" {
		http.Error(w, "ID and version parameters are required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var config manager_workflow.ConfigVersion
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	config.ID = id
	config.Version = version

	if err := h.repo.Update(&config); err != nil {
		log.Printf("Error updating config: %v", err)
		http.Error(w, "Failed to update config", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config)
}
