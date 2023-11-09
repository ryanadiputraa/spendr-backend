package domain

import (
	"context"
	"errors"
	"time"
)

type Expense struct {
	ID         string    `json:"id" db:"id"`
	CategoryID string    `json:"category_id" db:"category_id"`
	UserID     string    `json:"-" db:"user_id"`
	Expense    string    `json:"expense" db:"expense"`
	Amount     int       `json:"amount" db:"amount"`
	Date       time.Time `json:"date" db:"date"`
	CreatedAt  time.Time `json:"-" db:"created_at"`
	UpdatedAt  time.Time `json:"-" db:"updated_at"`
}

type ExpenseCategory struct {
	ID       string `json:"id" db:"id"`
	Category string `json:"category" db:"category"`
	Ico      string `json:"ico" db:"ico"`
	UserID   string `json:"-" db:"user_id"`
}

type ExpenseDTO struct {
	CategoryID string `json:"category_id" validate:"required"`
	Expense    string `json:"expense" validate:"required"`
	Amount     int    `json:"amount" validate:"required"`
	Date       string `json:"date" validate:"required,iso8601date"`
}

type ExpenseInfoDTO struct {
	ID          string    `json:"id"`
	Category    string    `json:"category"`
	CategoryIco string    `json:"category_ico" db:"category_ico"`
	Expense     string    `json:"expense"`
	Amount      int       `json:"amount"`
	Date        time.Time `json:"date"`
}

type ExpenseCategoryDTO struct {
	Category string `json:"category" validate:"required"`
	Ico      string `json:"ico" validate:"required,http_url"`
}

type ExpenseFilter struct {
	Size      int    `query:"size"`
	Page      int    `query:"page"`
	StartDate string `query:"start_date"`
	EndDate   string `query:"end_date"`
}

func NewExpense(id, userID string, dto ExpenseDTO) (*Expense, error) {
	date, err := time.Parse(time.RFC3339Nano, dto.Date)
	if err != nil {
		return nil, errors.New("invalid date format")
	}
	return &Expense{
		ID:         id,
		CategoryID: dto.CategoryID,
		UserID:     userID,
		Expense:    dto.Expense,
		Amount:     dto.Amount,
		Date:       date,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}, nil
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
	AddExpense(ctx context.Context, userID string, dto ExpenseDTO) (*Expense, error)
	ListExpense(ctx context.Context, userID string, filter ExpenseFilter) ([]ExpenseInfoDTO, error)
	DeleteExpense(ctx context.Context, userID, expenseID string) error
	AddExpenseCategory(ctx context.Context, userID string, dto ExpenseCategoryDTO) (*ExpenseCategory, error)
	ListExpenseCategory(ctx context.Context, userID string) ([]ExpenseCategory, error)
}

type ExpenseRepository interface {
	AddExpense(ctx context.Context, expense Expense) error
	ListExpense(ctx context.Context, userID string, filter ExpenseFilter) ([]ExpenseInfoDTO, error)
	DeleteExpense(ctx context.Context, userID, expenseID string) error
	AddExpenseCategory(ctx context.Context, category ExpenseCategory) error
	ListExpenseCategory(ctx context.Context, userID string) ([]ExpenseCategory, error)
}
