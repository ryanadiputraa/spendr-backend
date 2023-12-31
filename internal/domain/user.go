package domain

import (
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Currency string

const (
	IDR Currency = "idr"
	USD Currency = "usd"
	GBP Currency = "gbp"
	EUR Currency = "eur"
	YEN Currency = "yen"
)

// currencies is a map of supported currency
var Currencies map[string]bool = map[string]bool{
	"idr": true,
	"usd": true,
	"gbp": true,
	"eur": true,
	"yen": true,
}

type User struct {
	ID        string    `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"`
	FirstName string    `json:"first_name" db:"first_name"`
	LastName  string    `json:"last_name" db:"last_name"`
	Picture   string    `json:"picture" db:"picture"`
	Currency  string    `json:"currency" db:"currency"`
	CreatedAt time.Time `json:"-" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
}

type UserDTO struct {
	Email     string `json:"email" db:"email" validate:"required,email"`
	Password  string `json:"password" db:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" db:"first_name" validate:"required"`
	LastName  string `json:"last_name" db:"last_name" validate:"required"`
	Picture   string `json:"picture" db:"picture" validate:"http_url"`
	Currency  string `json:"currency" db:"currency" validate:"required"`
}

func NewUser(id, email, password, firstName, lastName, picture, currency string) (*User, error) {
	if _, ok := Currencies[currency]; !ok && currency != "" {
		return nil, errors.New("invalid currency")
	}

	return &User{
		ID:        id,
		Email:     email,
		Password:  password,
		FirstName: firstName,
		LastName:  lastName,
		Picture:   picture,
		Currency:  currency,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}, nil
}

func (u *User) HashPassword() error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	if err != nil {
		return err
	}
	u.Password = string(hashed)
	return nil
}

type UserService interface {
	GetUserData(ctx context.Context, userID string) (*User, error)
	ListSupportedCurrency(ctx context.Context) []string
}

type UserRepository interface {
	AddUser(ctx context.Context, user *User) error
	FindUserByID(ctx context.Context, userID string) (*User, error)
	FindUserByEmail(ctx context.Context, email string) (*User, error)
}
