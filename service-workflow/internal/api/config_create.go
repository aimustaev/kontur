package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aimustaev/service-workflow/internal/manager_workflow"
	"github.com/google/uuid"
)

type CreateConfigHandler struct {
	repo ConfigVersionRepository
}

func NewCreateConfigHandler(repo ConfigVersionRepository) *CreateConfigHandler {
	return &CreateConfigHandler{
		repo: repo,
	}
}

func (h *CreateConfigHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var config manager_workflow.ConfigVersion
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if config.Name == "" || config.Version == "" {
		http.Error(w, "Name and version are required", http.StatusBadRequest)
		return
	}

	config.ID = uuid.New()
	if err := h.repo.Create(&config); err != nil {
		log.Printf("Error creating config: %v", err)
		http.Error(w, "Failed to create config", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(config)
}
