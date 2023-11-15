package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	config "github.com/ryanadiputraa/spendr-backend/configs"
	"github.com/ryanadiputraa/spendr-backend/internal/domain"
)

type JWTService interface {
	GenerateJWTTokens(config *config.Config, claims domain.JWTClaims) (domain.JWTTokens, error)
	ParseJWTClaims(secret, tokenString string) (domain.JWTClaims, error)
}

type jwtService struct{}

func NewJWTService() JWTService {
	return &jwtService{}
}

func (j *jwtService) GenerateJWTTokens(config *config.Config, claims domain.JWTClaims) (domain.JWTTokens, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": claims.UserID,
		"exp":     time.Now().Add(config.ExpiresIn).Unix(),
	})
	tokenString, err := token.SignedString([]byte(config.Secret))
	if err != nil {
		return domain.JWTTokens{}, err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": claims.UserID,
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

func (j *jwtService) ParseJWTClaims(secret, tokenString string) (domain.JWTClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signin method %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return domain.JWTClaims{}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return domain.JWTClaims{}, err
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return domain.JWTClaims{}, errors.New("fail to parse jwt claims")
	}
	return domain.JWTClaims{
		UserID: userID,
	}, nil

}
