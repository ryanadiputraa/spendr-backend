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

func NewJWTService() JWTService {
	return &jwtService{}
}

func (j *jwtService) GenerateJWTTokens(config *config.Config, userID string) (domain.JWTTokens, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(config.ExpiresIn).Unix(),
	})
	tokenString, err := token.SignedString([]byte(config.Secret))
	if err != nil {
		return domain.JWTTokens{}, err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(config.RefreshExpiresIn).Unix(),
	})

	refreshTokenString, err := refreshToken.SignedString([]byte(config.RefreshSecret))
	if err != nil {
		return domain.JWTTokens{}, err
	}

	return domain.JWTTokens{
		AccessToken:  tokenString,
		ExpiresIn:    int(time.Now().Add(config.ExpiresIn).Unix()),
		RefreshToken: refreshTokenString,
	}, nil
}
