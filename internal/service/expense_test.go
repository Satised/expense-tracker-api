package service

import (
	"errors"
	"testing"

	"github.com/shopspring/decimal"

	"expense-tracker-api/internal/model"
)

type fakeRepo struct {
	added []model.Expense
}

func (f *fakeRepo) Add(e model.Expense) model.Expense {
	e.ID = len(f.added) + 1
	f.added = append(f.added, e)
	return e
}

func (f *fakeRepo) List() []model.Expense {
	return f.added
}

func TestCreate_OK(t *testing.T) {
	repo := &fakeRepo{}
	svc := NewExpenseService(repo)

	got, err := svc.Create(CreateExpenseInput{
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

			_, err := svc.Create(CreateExpenseInput{
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
