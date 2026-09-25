package service

import (
	"context"
	"errors"
	"testing"

	"github.com/shopspring/decimal"

	"expense-tracker-api/internal/model"
)

type fakeRepo struct {
	added []model.Expense
}

func (f *fakeRepo) Add(ctx context.Context, e model.Expense) (model.Expense, error) {
	e.ID = len(f.added) + 1
	f.added = append(f.added, e)
	return e, nil
}

func (f *fakeRepo) List(ctx context.Context) ([]model.Expense, error) {
	return f.added, nil
}

func TestCreate_OK(t *testing.T) {
	repo := &fakeRepo{}
	svc := NewExpenseService(repo)

	got, err := svc.Create(context.Background(), CreateExpenseInput{
		Amount:   decimal.NewFromInt(1500),
		Category: "food",
	})

	if err != nil {
		t.Fatalf("Create() returned error: %v", err)
	}

	if got.ID != 1 {
		t.Errorf("Create() returned ID = %d, want 1", got.ID)
	}
}

func TestCreate_Errors(t *testing.T) {
	tests := []struct {
		name     string
		amount   decimal.Decimal
		category string
		wantErr  error
	}{
		{
			name:     "пустая категория",
			amount:   decimal.NewFromInt(1500),
			category: "",
			wantErr:  ErrCategoryRequired,
		},
		{
			name:     "категория из пробелов",
			amount:   decimal.NewFromInt(1500),
			category: "   ",
			wantErr:  ErrCategoryRequired,
		},
		{
			name:     "отрицательная сумма",
			amount:   decimal.NewFromInt(-5),
			category: "food",
			wantErr:  ErrAmountNotPositive,
		},
		{
			name:     "нулевая сумма",
			amount:   decimal.NewFromInt(0),
			category: "food",
			wantErr:  ErrAmountNotPositive,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := &fakeRepo{}
			svc := NewExpenseService(repo)

			_, err := svc.Create(context.Background(), CreateExpenseInput{
				Amount:   tc.amount,
				Category: tc.category,
			})

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Create() returned error = %v, want %v", err, tc.wantErr)
			}

			if len(repo.added) != 0 {
				t.Errorf("Create() added expense to repo, want no expenses added")
			}
		})
	}
}
