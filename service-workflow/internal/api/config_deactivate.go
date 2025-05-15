package api

import (
	"log"
	"net/http"

	"github.com/aimustaev/service-workflow/internal/manager_workflow"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type DeactivateConfigHandler struct {
	repo manager_workflow.ConfigVersionRepository
}

func NewDeactivateConfigHandler(repo manager_workflow.ConfigVersionRepository) *DeactivateConfigHandler {
	return &DeactivateConfigHandler{
		repo: repo,
	}
}

func (h *DeactivateConfigHandler) Handle(w http.ResponseWriter, r *http.Request) {
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

	if err := h.repo.Deactivate(id, version); err != nil {
		log.Printf("Error deactivating config: %v", err)
		http.Error(w, "Failed to deactivate config", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
