package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ryanadiputraa/spendr-backend/config"
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

		tokens, err := h.jwtService.GenerateJWTTokens(h.config, user.ID)
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
