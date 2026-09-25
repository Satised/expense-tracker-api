package repository

import (
	"context"
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

func (s *ExpenseStore) Add(ctx context.Context, e model.Expense) (model.Expense, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	e.ID = s.nextID
	s.nextID++
	e.CreatedAt = time.Now()

	s.items = append(s.items, e)
	return e, nil
}

func (s *ExpenseStore) List(ctx context.Context) ([]model.Expense, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := make([]model.Expense, len(s.items))
	copy(out, s.items)
	return out, nil
}
