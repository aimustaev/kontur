package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type GetVersionConfigHandler struct {
	repo ConfigVersionRepository
}

func NewGetVersionConfigHandler(repo ConfigVersionRepository) *GetVersionConfigHandler {
	return &GetVersionConfigHandler{
		repo: repo,
	}
}

func (h *GetVersionConfigHandler) Handle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	version := vars["version"]
	if name == "" || version == "" {
		http.Error(w, "Name and version parameters are required", http.StatusBadRequest)
		return
	}

	config, err := h.repo.GetByVersion(name, version)
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
