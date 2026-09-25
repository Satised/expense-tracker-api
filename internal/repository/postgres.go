package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"expense-tracker-api/internal/model"
)

type PostgresStore struct {
	pool *pgxpool.Pool
}

func NewPostgresStore(pool *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{pool: pool}
}

func (s *PostgresStore) Add(ctx context.Context, e model.Expense) (model.Expense, error) {
	const query = `
		INSERT INTO expenses (amount, category, note)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`

	err := s.pool.QueryRow(ctx, query, e.Amount, e.Category, e.Note).
		Scan(&e.ID, &e.CreatedAt)
	if err != nil {
		return model.Expense{}, fmt.Errorf("insert expense: %w", err)
	}

	return e, nil
}

func (s *PostgresStore) List(ctx context.Context) ([]model.Expense, error) {
	const query = `
		SELECT id, amount, category, note, created_at
		FROM expenses
		ORDER BY created_at DESC`

	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query expenses: %w", err)
	}
	defer rows.Close()

	expenses := []model.Expense{}

	for rows.Next() {
		var e model.Expense

		err := rows.Scan(&e.ID, &e.Amount, &e.Category, &e.Note, &e.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan expense: %w", err)
		}

		expenses = append(expenses, e)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read expenses: %w", err)
	}

	return expenses, nil
}
