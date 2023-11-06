package domain

import (
	"context"
	"errors"

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
var currencies map[string]bool = map[string]bool{
	"idr": true,
	"usd": true,
	"gbp": true,
	"eur": true,
	"yen": true,
}

type User struct {
	ID        string `json:"id" db:"id"`
	Email     string `json:"email" db:"email"`
	Password  string `json:"-" db:"password"`
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`
	Picture   string `json:"picture" db:"picture"`
	Currency  string `json:"currency" db:"currency"`
}

type UserDTO struct {
	Email     string `json:"email" db:"email"`
	Password  string `json:"-" db:"password"`
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`
	Picture   string `json:"picture" db:"picture"`
	Currency  string `json:"currency" db:"currency"`
}

func NewUser(id, email, password, firstName, lastName, picture, currency string) (*User, error) {
	if _, ok := currencies[currency]; !ok {
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
	RegisterUser(ctx context.Context, dto UserDTO) (*User, error)
}

type UserRepository interface {
	AddUser(ctx context.Context, user User) (*User, error)
}
