package domain

import "context"

type SigninDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type JWTTokens struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type AuthService interface {
	Signup(ctx context.Context, dto UserDTO) (*User, error)
	Signin(ctx context.Context, email, password string) (*User, error)
}
