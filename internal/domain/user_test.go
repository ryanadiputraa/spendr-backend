package domain

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

const (
	DummyImg = "https://placehold.co/120x80"
)

func TestNewUser(t *testing.T) {
	now := time.Now().UTC()

	user := User{
		ID:        uuid.NewString(),
		Email:     "john@mail.com",
		Password:  "hashedpassword",
		FirstName: "John",
		LastName:  "Doe",
		Picture:   DummyImg,
		Currency:  "idr",
		CreatedAt: now,
		UpdatedAt: now,
	}
	invalidUserCurrency := User{
		ID:        uuid.NewString(),
		Email:     "jane@mail.com",
		Password:  "hashedpassword",
		FirstName: "Jane",
		LastName:  "Doe",
		Picture:   DummyImg,
		Currency:  "won",
		CreatedAt: now,
		UpdatedAt: now,
	}

	cases := []struct {
		name     string
		payload  User
		expected *User
		err      error
	}{
		{
			name:     "should return a new user",
			expected: &user,
			payload:  user,
			err:      nil,
		},
		{
			name:     "should return error invalid currency",
			expected: nil,
			payload:  invalidUserCurrency,
			err:      errors.New("invalid currency"),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			u, err := NewUser(
				c.payload.ID,
				c.payload.Email,
				c.payload.Password,
				c.payload.FirstName,
				c.payload.LastName,
				c.payload.Picture,
				c.payload.Currency,
			)
			assert.Equal(t, c.err, err)
			if err != nil {
				assert.Empty(t, u)
				return
			}

			assert.NotEmpty(t, u.ID)
			assert.Equal(t, c.expected.Email, u.Email)
			assert.Equal(t, c.expected.Password, u.Password)
			assert.Equal(t, c.expected.FirstName, u.FirstName)
			assert.Equal(t, c.expected.LastName, u.LastName)
			assert.Equal(t, c.expected.Picture, u.Picture)
			assert.Equal(t, c.expected.Currency, u.Currency)
			assert.WithinDuration(t, user.CreatedAt, u.CreatedAt, 500*time.Millisecond)
			assert.WithinDuration(t, user.UpdatedAt, u.UpdatedAt, 500*time.Millisecond)
		})
	}
}

func TestHashPassword(t *testing.T) {
	password := "plainpassword"
	user := User{
		ID:        uuid.NewString(),
		Email:     "john@mail.com",
		Password:  password,
		FirstName: "John",
		LastName:  "Doe",
		Picture:   DummyImg,
		Currency:  "idr",
	}

	cases := []struct {
		name string
		user User
		err  error
	}{
		{
			name: "should successfully hash user password",
			user: user,
			err:  nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.user.HashPassword()
			assert.Equal(t, c.err, err)

			if err != nil {
				err = bcrypt.CompareHashAndPassword([]byte(c.user.Password), []byte(password))
				assert.Nil(t, err)
			}
		})
	}
}
