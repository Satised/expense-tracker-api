package repository

import (
	"sync"
	"time"

	"expense-tracker-api/internal/model"
)

type ExpenseStore struct {
	mu     sync.Mutex
	items  []model.Expense
	nextID int
}

func NewExpenseStore() *ExpenseStore {
	return &ExpenseStore{nextID: 1}
}

func (s *ExpenseStore) Add(e model.Expense) model.Expense {
	s.mu.Lock()
	defer s.mu.Unlock()

	e.ID = s.nextID
	s.nextID++
	e.CreatedAt = time.Now()

	s.items = append(s.items, e)
	return e
}

func (s *ExpenseStore) List() []model.Expense {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := make([]model.Expense, len(s.items))
	copy(out, s.items)
	return out
}
