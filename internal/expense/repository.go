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
