package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type DeactivateConfigHandler struct {
	repo ConfigVersionRepository
}

func NewDeactivateConfigHandler(repo ConfigVersionRepository) *DeactivateConfigHandler {
	return &DeactivateConfigHandler{
		repo: repo,
	}
}

func (h *DeactivateConfigHandler) Handle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	version := vars["version"]
	if name == "" || version == "" {
		http.Error(w, "Name and version parameters are required", http.StatusBadRequest)
		return
	}

	if err := h.repo.Deactivate(name, version); err != nil {
		log.Printf("Error deactivating config: %v", err)
		http.Error(w, "Failed to deactivate config", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
