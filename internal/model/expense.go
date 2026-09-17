package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Expense struct {
	ID        int             `json:"id"`
	Amount    decimal.Decimal `json:"amount"`
	Category  string          `json:"category"`
	Note      string          `json:"note"`
	CreatedAt time.Time       `json:"created_at"`
}
