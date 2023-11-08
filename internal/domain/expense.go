package domain

import "time"

type Expense struct {
	ID         string    `json:"id" db:"id"`
	CategoryID string    `json:"category_id" db:"category_id"`
	UserID     string    `json:"user_id" db:"user_id"`
	Expense    string    `json:"expense" db:"expense"`
	Amount     int       `json:"amount" db:"amount"`
	Date       time.Time `json:"date" db:"date"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type ExpenseCategory struct {
	ID       string `json:"id" db:"id"`
	Category string `json:"category" db:"category"`
	Ico      string `json:"ico" db:"ico"`
	UserID   string `json:"user_id" db:"user_id"`
}
