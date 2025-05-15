package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aimustaev/service-workflow/internal/manager_workflow"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type GetVersionConfigHandler struct {
	repo manager_workflow.ConfigVersionRepository
}

func NewGetVersionConfigHandler(repo manager_workflow.ConfigVersionRepository) *GetVersionConfigHandler {
	return &GetVersionConfigHandler{
		repo: repo,
	}
}

func (h *GetVersionConfigHandler) Handle(w http.ResponseWriter, r *http.Request) {
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

	config, err := h.repo.GetByVersion(id, version)
	if err != nil {
		log.Printf("Error getting config version: %v", err)
		http.Error(w, "Failed to get config version", http.StatusInternalServerError)
		return
	}
	if config == nil {
		http.Error(w, "Config not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config)
}
