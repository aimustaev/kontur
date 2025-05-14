package api

import (
	"net/http"
)

type HealthHandler struct{}

func (h *HealthHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
