package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	rpc "github.com/aimustaev/service-gateway/internal/tickets"
)

type Handler struct {
	ticketsClient *rpc.Client
}

func NewHandler(ticketsClient *rpc.Client) *Handler {
	return &Handler{
		ticketsClient: ticketsClient,
	}
}

func (h *Handler) GetAllTickets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response, err := h.ticketsClient.GetAllTickets(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetTicketMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	ticketID := vars["id"]
	if ticketID == "" {
		http.Error(w, "Invalid ticket ID", http.StatusBadRequest)
		return
	}

	response, err := h.ticketsClient.GetTicketMessages(r.Context(), ticketID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
