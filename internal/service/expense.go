package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/shopspring/decimal"

	"expense-tracker-api/internal/model"
)

// --- ошибки бизнес-правил ---

var (
	ErrCategoryRequired  = errors.New("category is required")
	ErrAmountNotPositive = errors.New("amount must be greater than zero")
)

// --- что сервис требует от хранилища ---

type ExpenseRepository interface {
	Add(ctx context.Context, e model.Expense) (model.Expense, error)
	List(ctx context.Context) ([]model.Expense, error)
}

// --- сам сервис ---

type ExpenseService struct {
	repo ExpenseRepository
}

func NewExpenseService(repo ExpenseRepository) *ExpenseService {
	return &ExpenseService{repo: repo}
}

// --- то, что приходит на создание расхода ---

type CreateExpenseInput struct {
	Amount   decimal.Decimal
	Category string
	Note     string
}

func (s *ExpenseService) Create(ctx context.Context, in CreateExpenseInput) (model.Expense, error) {
	category := strings.TrimSpace(in.Category)

	if category == "" {
		return model.Expense{}, ErrCategoryRequired
	}

	if !in.Amount.IsPositive() {
		return model.Expense{}, ErrAmountNotPositive
	}

	created, err := s.repo.Add(ctx, model.Expense{
		Amount:   in.Amount,
		Category: category,
		Note:     strings.TrimSpace(in.Note),
	})
	if err != nil {
		return model.Expense{}, fmt.Errorf("add expense: %w", err)
	}

	return created, nil
}

func (s *ExpenseService) List(ctx context.Context) ([]model.Expense, error) {
	expenses, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list expenses: %w", err)
	}

	return expenses, nil
}
