package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ryanadiputraa/spendr-backend/config"
	"github.com/ryanadiputraa/spendr-backend/internal/domain"
)

type JWTService interface {
	GenerateJWTTokens(config *config.Config, userID string) (domain.JWTTokens, error)
}

type jwtService struct{}

func (j *jwtService) GenerateJWTTokens(config *config.Config, userID string) (domain.JWTTokens, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(config.ExpiresIn).Unix(),
	})
	tokenString, err := token.SignedString(config.JWT.Secret)
	if err != nil {
		return domain.JWTTokens{}, nil
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(config.RefreshExpiresIn).Unix(),
	})

	refreshTokenString, err := refreshToken.SignedString(config.JWT.RefreshSecret)
	if err != nil {
		return domain.JWTTokens{}, nil
	}
	return domain.JWTTokens{
		AccessToken:  tokenString,
		ExpiresIn:    int(time.Now().Add(config.ExpiresIn).Unix()),
		RefreshToken: refreshTokenString,
	}, nil
}
