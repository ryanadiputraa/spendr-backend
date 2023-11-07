package user

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ryanadiputraa/spendr-backend/internal/domain"
	"github.com/ryanadiputraa/spendr-backend/pkg/validator"
)

type handler struct {
	validator validator.Validator
	service   domain.UserService
}

func NewHandler(group *echo.Group, validator validator.Validator, service domain.UserService) {
	h := handler{
		validator: validator,
		service:   service,
	}

	group.POST("/register", h.Signup())
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
			code, resp := domain.MapServiceErrHTTPResponse(err)
			return c.JSON(code, resp)
		}

		return c.JSON(http.StatusOK, map[string]any{
			"data": user,
		})
	}
}
