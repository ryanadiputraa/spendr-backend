package user

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/ryanadiputraa/spendr-backend/internal/domain"
)

type repository struct {
	DB *sqlx.DB
}

func NewRepository(DB *sqlx.DB) domain.UserRepository {
	return &repository{
		DB: DB,
	}
}

func (r *repository) AddUser(ctx context.Context, user *domain.User) error {
	q := `INSERT INTO users (id, email, password, first_name, last_name, picture, currency, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.DB.ExecContext(
		ctx,
		q,
		user.ID,
		user.Email,
		user.Password,
		user.FirstName,
		user.LastName,
		user.Picture,
		user.Currency,
		user.CreatedAt,
		user.UpdatedAt,
	)
	return err
}

func (r *repository) FindUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	q := `SELECT id, email, password, first_name, last_name, picture, currency FROM users WHERE email = $1`

	var user domain.User
	err := r.DB.Get(&user, q, email)

	return &user, err
}
