package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/ryanadiputraa/spendr-backend/config"
	"github.com/ryanadiputraa/spendr-backend/pkg/jwt"
)

type AuthMiddleware struct {
	config *config.Config
	jwt    jwt.JWTService
}

func NewAuthMiddleware(config *config.Config, jwt jwt.JWTService) *AuthMiddleware {
	return &AuthMiddleware{
		config: config,
		jwt:    jwt,
	}
}

func (m *AuthMiddleware) ParseJWTClaims(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		header := c.Request().Header
		auth := header["Authorization"]
		if len(auth) == 0 {
			return c.JSON(http.StatusUnauthorized, map[string]any{
				"error": "missing Authorization header",
			})
		}

		authToken := auth[0]
		tokens := strings.Split(authToken, " ")
		if len(tokens) < 2 || tokens[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, map[string]any{
				"error": "invalid token, expecting Bearer token header",
			})
		}

		claims, err := m.jwt.ParseJWTClaims(m.config.JWT.Secret, tokens[1])
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]any{
				"error": "fail to parse jwt tokens, " + err.Error(),
			})
		}
		c.Set("user_id", claims.UserID)

		return next(c)
	}
}
