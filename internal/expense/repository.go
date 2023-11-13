package expense

import (
	"context"
	"database/sql"

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

func (r *repository) ListExpense(ctx context.Context, userID string, filter domain.ExpenseFilter) ([]domain.ExpenseInfoDTO, error) {
	q := `SELECT expenses.id AS id, c.category AS category, c.ico AS category_ico, expense, amount, date
    FROM expenses LEFT JOIN expense_categories AS c ON c.id = expenses.category_id
    WHERE expenses.user_id = $1 AND date BETWEEN $4 AND $5
    ORDER BY date DESC LIMIT $2 OFFSET $3`

	var expenses []domain.ExpenseInfoDTO
	offset := (filter.Page - 1) * filter.Size
	err := r.DB.Select(&expenses, q, userID, filter.Size, offset, filter.StartDate, filter.EndDate)

	return expenses, err
}

func (r *repository) DeleteExpense(ctx context.Context, userID, expenseID string) error {
	q := `DELETE FROM expenses WHERE id = $1 AND user_id = $2`
	res, err := r.DB.Exec(q, expenseID, userID)
	c, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if c < 1 {
		return sql.ErrNoRows
	}

	return err
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

func (r *repository) DeleteExpenseCategory(ctx context.Context, userID, categoryID string) error {
	updateExpenseCategories := `UPDATE expenses SET category_id = NULL WHERE category_id = $1 AND user_id = $2`
	deleteCategories := `DELETE FROM expense_categories WHERE id = $1 AND user_id = $2`

	tx, err := r.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, updateExpenseCategories, categoryID, userID)
	if err != nil {
		tx.Rollback()
		return sql.ErrNoRows
	}

	res, err := tx.ExecContext(ctx, deleteCategories, categoryID, userID)
	if err != nil {
		tx.Rollback()
		return sql.ErrNoRows
	}
	c, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if c < 1 {
		return sql.ErrNoRows
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
