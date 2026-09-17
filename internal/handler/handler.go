package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"expense-tracker-api/internal/service"
)

type Handler struct {
	expenses *service.ExpenseService
}

func New(expenses *service.ExpenseService) *Handler {
	return &Handler{expenses: expenses}
}

// --- формы ответов ---

type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// --- помощники ---

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Println("failed to write json:", err)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{Error: message})
}

// --- служебная ручка ---

func (h *Handler) Healthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, HealthResponse{Status: "ok", Version: "0.1.0"})
}
