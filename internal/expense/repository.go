package expense

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/ryanadiputraa/spendr-backend/internal/domain"
)

type repository struct {
	DB *sqlx.DB
}

func NewRepository(DB *sqlx.DB) domain.ExpenseRepository {
	return &repository{
		DB: DB,
	}
}

func (r *repository) AddExpense(ctx context.Context, expense domain.Expense) error {
	q := `INSERT INTO expenses (id, category_id, user_id, expense, amount, date, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := r.DB.Exec(
		q,
		expense.ID,
		expense.CategoryID,
		expense.UserID,
		expense.Expense,
		expense.Amount,
		expense.Date,
		expense.CreatedAt,
		expense.UpdatedAt,
	)

	return err
}

func (r *repository) ListLatestExpense(ctx context.Context, userID string, limit int) ([]domain.ExpenseInfoDTO, error) {
	q := `SELECT expenses.id AS id, c.category AS category, c.ico AS category_ico, expense, amount, date
    FROM expenses LEFT JOIN expense_categories AS c ON c.id = expenses.category_id
    WHERE expenses.user_id = $1 ORDER BY date DESC LIMIT $2`

	var expenses []domain.ExpenseInfoDTO
	err := r.DB.Select(&expenses, q, userID, limit)

	return expenses, err
}

func (r *repository) AddExpenseCategory(ctx context.Context, category domain.ExpenseCategory) error {
	q := `INSERT INTO expense_categories (id, category, ico, user_id) VALUES ($1, $2, $3, $4)`
	_, err := r.DB.Exec(q, category.ID, category.Category, category.Ico, category.UserID)
	return err
}

func (r *repository) ListExpenseCategory(ctx context.Context, userID string) ([]domain.ExpenseCategory, error) {
	q := `SELECT id, category, ico FROM expense_categories WHERE user_id = $1 ORDER BY category ASC`

	var categories []domain.ExpenseCategory
	err := r.DB.Select(&categories, q, userID)

	return categories, err
}
