package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/shopspring/decimal"

	"expense-tracker-api/internal/service"
)

func (h *Handler) CreateExpense(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Amount   decimal.Decimal `json:"amount"`
		Category string          `json:"category"`
		Note     string          `json:"note"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid json body")
		return
	}

	created, err := h.expenses.Create(r.Context(), service.CreateExpenseInput{
		Amount:   req.Amount,
		Category: req.Category,
		Note:     req.Note,
	})

	if err != nil {
		if errors.Is(err, service.ErrCategoryRequired) ||
			errors.Is(err, service.ErrAmountNotPositive) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		log.Println("create expense failed:", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

func (h *Handler) ListExpenses(w http.ResponseWriter, r *http.Request) {
	expenses, err := h.expenses.List(r.Context())
	if err != nil {
		log.Println("list expenses failed:", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusOK, expenses)
}
