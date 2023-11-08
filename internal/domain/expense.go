package domain

import (
	"context"
	"time"
)

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
	UserID   string `json:"-" db:"user_id"`
}

type ExpenseCategoryDTO struct {
	Category string `json:"category" db:"category" validate:"required"`
	Ico      string `json:"ico" db:"ico" validate:"required,http_url"`
}

func NewExpenseCategory(id, category, ico, userID string) *ExpenseCategory {
	return &ExpenseCategory{
		ID:       id,
		Category: category,
		Ico:      ico,
		UserID:   userID,
	}
}

type ExpenseService interface {
	AddExpenseCategory(ctx context.Context, userID string, dto ExpenseCategoryDTO) (*ExpenseCategory, error)
}

type ExpenseRepository interface {
	AddExpenseCategory(ctx context.Context, category ExpenseCategory) error
}
