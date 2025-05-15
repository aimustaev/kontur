package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"

	"github.com/aimustaev/service-workflow/internal/manager_workflow"
)

// CreateConfigRequest представляет запрос на создание конфигурации
type CreateConfigRequest struct {
	Name      string          `json:"name"`
	Version   string          `json:"version"`
	Content   json.RawMessage `json:"content"` // Может быть как JSON объектом, так и строкой
	Schema    json.RawMessage `json:"schema,omitempty"`
	CreatedBy string          `json:"created_by"`
	IsActive  bool            `json:"is_active"`
}

type CreateConfigHandler struct {
	repo manager_workflow.ConfigVersionRepository
}

func NewCreateConfigHandler(repo manager_workflow.ConfigVersionRepository) *CreateConfigHandler {
	return &CreateConfigHandler{
		repo: repo,
	}
}

func (h *CreateConfigHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var req CreateConfigRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Version == "" {
		http.Error(w, "Name and version are required", http.StatusBadRequest)
		return
	}

	// Проверяем, что content является валидным JSON
	var contentBytes []byte
	if len(req.Content) > 0 {
		// Если content - строка, то она должна быть в кавычках
		if req.Content[0] == '"' && req.Content[len(req.Content)-1] == '"' {
			// Убираем кавычки и экранирование
			var contentStr string
			if err := json.Unmarshal(req.Content, &contentStr); err != nil {
				http.Error(w, "Invalid content format: string must be properly escaped", http.StatusBadRequest)
				return
			}
			contentBytes = []byte(contentStr)
		} else {
			// Проверяем, что это валидный JSON
			var contentObj interface{}
			if err := json.Unmarshal(req.Content, &contentObj); err != nil {
				http.Error(w, "Invalid content format: must be a valid JSON object or string", http.StatusBadRequest)
				return
			}
			contentBytes = req.Content
		}
	}

	config := &manager_workflow.ConfigVersion{
		ID:        uuid.New(),
		Name:      req.Name,
		Version:   req.Version,
		Content:   contentBytes,
		Schema:    &req.Schema,
		CreatedBy: req.CreatedBy,
		IsActive:  req.IsActive,
	}

	if err := h.repo.Create(config); err != nil {
		log.Printf("Error creating config: %v", err)
		http.Error(w, "Failed to create config", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(config)
}
