package user

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ryanadiputraa/spendr-backend/internal/domain"
	"github.com/ryanadiputraa/spendr-backend/internal/middleware"
	"github.com/ryanadiputraa/spendr-backend/pkg/httpres"
)

type handler struct {
	service        domain.UserService
	authMiddleware middleware.AuthMiddleware
}

func NewHandler(group *echo.Group, service domain.UserService, authMiddleware middleware.AuthMiddleware) {
	h := &handler{
		service: service,
	}

	group.GET("", h.GetUserData(), authMiddleware.ParseJWTClaims)
	group.GET("/currencies", h.ListSupportedCurrency())
}

func (h *handler) GetUserData() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := c.Get("user_id").(string)

		user, err := h.service.GetUserData(c.Request().Context(), userID)
		if err != nil {
			code, resp := httpres.MapServiceErrHTTPResponse(err)
			return c.JSON(code, resp)
		}

		return c.JSON(http.StatusOK, map[string]any{
			"data": user,
		})
	}
}

func (h *handler) ListSupportedCurrency() echo.HandlerFunc {
	return func(c echo.Context) error {
		currencies := h.service.ListSupportedCurrency(c.Request().Context())

		return c.JSON(http.StatusOK, map[string]any{
			"data": currencies,
		})
	}
}
