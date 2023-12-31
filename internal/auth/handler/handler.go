package auth

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	config "github.com/ryanadiputraa/spendr-backend/configs"
	"github.com/ryanadiputraa/spendr-backend/internal/domain"
	"github.com/ryanadiputraa/spendr-backend/pkg/httpres"
	"github.com/ryanadiputraa/spendr-backend/pkg/jwt"
	"github.com/ryanadiputraa/spendr-backend/pkg/logger"
	"github.com/ryanadiputraa/spendr-backend/pkg/validator"
)

type handler struct {
	config     *config.Config
	log        logger.Logger
	validator  validator.Validator
	service    domain.AuthService
	jwtService jwt.JWTService
}

func NewHandler(group *echo.Group, config *config.Config, log logger.Logger, validator validator.Validator, service domain.AuthService, jwtService jwt.JWTService) {
	h := handler{
		config:     config,
		log:        log,
		validator:  validator,
		service:    service,
		jwtService: jwtService,
	}

	group.POST("/register", h.Signup())
	group.POST("/login", h.Signin())
	group.GET("/refresh_token", h.RefreshJWTTokens())
}

func (h *handler) Signup() echo.HandlerFunc {
	return func(c echo.Context) error {
		var dto domain.UserDTO
		if err := c.Bind(&dto); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"error": err.Error(),
			})
		}

		err, errDetails := h.validator.Validate(dto)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"error":  err.Error(),
				"errors": errDetails,
			})
		}

		user, err := h.service.Signup(c.Request().Context(), dto)
		if err != nil {
			code, resp := httpres.MapServiceErrHTTPResponse(err)
			return c.JSON(code, resp)
		}

		return c.JSON(http.StatusOK, map[string]any{
			"data": user,
		})
	}
}

func (h *handler) Signin() echo.HandlerFunc {
	return func(c echo.Context) error {
		var dto domain.SigninDTO
		if err := c.Bind(&dto); err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]any{
				"error": err.Error(),
			})
		}

		user, err := h.service.Signin(c.Request().Context(), dto.Email, dto.Password)
		if err != nil {
			code, resp := httpres.MapServiceErrHTTPResponse(err)
			return c.JSON(code, resp)
		}

		tokens, err := h.jwtService.GenerateJWTTokens(h.config, domain.JWTClaims{
			UserID: user.ID,
		})
		if err != nil {
			h.log.Error("auth handler: ", err)
			return c.JSON(http.StatusInternalServerError, map[string]any{
				"error": "internal server error",
			})
		}

		return c.JSON(http.StatusOK, map[string]any{
			"data": tokens,
		})
	}
}

func (h *handler) RefreshJWTTokens() echo.HandlerFunc {
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

		claims, err := h.jwtService.ParseJWTClaims(h.config.JWT.RefreshSecret, tokens[1])
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]any{
				"error": "fail to parse jwt claims" + err.Error(),
			})
		}

		refreshedTokens, err := h.jwtService.GenerateJWTTokens(h.config, domain.JWTClaims{
			UserID: claims.UserID,
		})
		if err != nil {
			h.log.Error("auth handler: ", err)
			return c.JSON(http.StatusInternalServerError, map[string]any{
				"error": "internal server error",
			})
		}

		return c.JSON(http.StatusOK, map[string]any{
			"data": refreshedTokens,
		})
	}
}
