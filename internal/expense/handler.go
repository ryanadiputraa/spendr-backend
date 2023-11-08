package expense

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ryanadiputraa/spendr-backend/internal/domain"
	"github.com/ryanadiputraa/spendr-backend/internal/middleware"
	"github.com/ryanadiputraa/spendr-backend/pkg/httpres"
	"github.com/ryanadiputraa/spendr-backend/pkg/logger"
	"github.com/ryanadiputraa/spendr-backend/pkg/validator"
)

type handler struct {
	log            logger.Logger
	validator      validator.Validator
	service        domain.ExpenseService
	authMiddleware middleware.AuthMiddleware
}

func NewHandler(group *echo.Group, log logger.Logger, validator validator.Validator, service domain.ExpenseService, authMiddleware middleware.AuthMiddleware) {
	h := &handler{
		log:       log,
		validator: validator,
		service:   service,
	}

	group.POST("/categories", h.AddExpenseCategory(), authMiddleware.ParseJWTClaims)
}

func (h *handler) AddExpenseCategory() echo.HandlerFunc {
	return func(c echo.Context) error {
		var dto domain.ExpenseCategoryDTO
		if err := c.Bind(&dto); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"error": err.Error(),
			})
		}

		err, errors := h.validator.Validate(dto)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"error":  err.Error(),
				"errros": errors,
			})
		}

		userID := c.Get("user_id").(string)
		category, err := h.service.AddExpenseCategory(c.Request().Context(), userID, dto)
		if err != nil {
			code, resp := httpres.MapServiceErrHTTPResponse(err)
			return c.JSON(code, resp)
		}

		return c.JSON(http.StatusCreated, map[string]any{
			"data": category,
		})
	}
}
