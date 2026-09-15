package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/shopspring/decimal"
)

// --- модели ---

type Expense struct {
	ID        int             `json:"id"`
	Amount    decimal.Decimal `json:"amount"`
	Category  string          `json:"category"`
	Note      string          `json:"note"`
	CreatedAt time.Time       `json:"created_at"`
}

type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// --- хранилище в памяти ---

type ExpenseStore struct {
	mu     sync.Mutex
	items  []Expense
	nextID int
}

func NewExpenseStore() *ExpenseStore {
	return &ExpenseStore{nextID: 1}
}

func (s *ExpenseStore) Add(e Expense) Expense {
	s.mu.Lock()
	defer s.mu.Unlock()

	e.ID = s.nextID
	s.nextID++
	e.CreatedAt = time.Now()

	s.items = append(s.items, e)
	return e
}

func (s *ExpenseStore) List() []Expense {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := make([]Expense, len(s.items))
	copy(out, s.items)
	return out
}

// --- помощники для ответов ---

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

// --- хендлеры ---

type API struct {
	expenses *ExpenseStore
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, HealthResponse{Status: "ok", Version: "0.1.0"})
}

func (a *API) createExpense(w http.ResponseWriter, r *http.Request) {
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

	req.Category = strings.TrimSpace(req.Category)

	if req.Category == "" {
		writeError(w, http.StatusBadRequest, "category is required")
		return
	}

	if !req.Amount.IsPositive() {
		writeError(w, http.StatusBadRequest, "amount must be greater than zero")
		return
	}

	created := a.expenses.Add(Expense{
		Amount:   req.Amount,
		Category: req.Category,
		Note:     strings.TrimSpace(req.Note),
	})

	writeJSON(w, http.StatusCreated, created)
}

func (a *API) listExpenses(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, a.expenses.List())
}

// --- точка входа ---

func main() {
	api := &API{expenses: NewExpenseStore()}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", healthzHandler)
	mux.HandleFunc("POST /expenses", api.createExpense)
	mux.HandleFunc("GET /expenses", api.listExpenses)

	log.Println("сервер слушает на http://localhost:8080")

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
