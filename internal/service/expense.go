package service

import (
	"errors"
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
	Add(e model.Expense) model.Expense
	List() []model.Expense
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

func (s *ExpenseService) Create(in CreateExpenseInput) (model.Expense, error) {
	category := strings.TrimSpace(in.Category)

	if category == "" {
		return model.Expense{}, ErrCategoryRequired
	}

	if !in.Amount.IsPositive() {
		return model.Expense{}, ErrAmountNotPositive
	}

	created := s.repo.Add(model.Expense{
		Amount:   in.Amount,
		Category: category,
		Note:     strings.TrimSpace(in.Note),
	})

	return created, nil
}

func (s *ExpenseService) List() []model.Expense {
	return s.repo.List()
}
