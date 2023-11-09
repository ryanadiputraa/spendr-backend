package expense

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ryanadiputraa/spendr-backend/internal/domain"
	"github.com/ryanadiputraa/spendr-backend/internal/middleware"
	"github.com/ryanadiputraa/spendr-backend/pkg/httpres"
	"github.com/ryanadiputraa/spendr-backend/pkg/validator"
)

type handler struct {
	validator      validator.Validator
	service        domain.ExpenseService
	authMiddleware middleware.AuthMiddleware
}

func NewHandler(group *echo.Group, validator validator.Validator, service domain.ExpenseService, authMiddleware middleware.AuthMiddleware) {
	h := &handler{
		validator: validator,
		service:   service,
	}

	group.POST("", h.AddExpense(), authMiddleware.ParseJWTClaims)
	group.GET("", h.ListExpense(), authMiddleware.ParseJWTClaims)
	group.GET("/categories", h.ListExpenseCategory(), authMiddleware.ParseJWTClaims)
	group.POST("/categories", h.AddExpenseCategory(), authMiddleware.ParseJWTClaims)
}

func (h *handler) AddExpense() echo.HandlerFunc {
	return func(c echo.Context) error {
		var dto domain.ExpenseDTO
		if err := c.Bind(&dto); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"error": err.Error(),
			})
		}

		err, errors := h.validator.Validate(dto)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"error":  err.Error(),
				"errors": errors,
			})
		}

		userID := c.Get("user_id").(string)
		expense, err := h.service.AddExpense(c.Request().Context(), userID, dto)
		if err != nil {
			code, resp := httpres.MapServiceErrHTTPResponse(err)
			return c.JSON(code, resp)
		}

		return c.JSON(http.StatusCreated, map[string]any{
			"data": expense,
		})
	}
}

func (h *handler) ListExpense() echo.HandlerFunc {
	return func(c echo.Context) error {
		var params domain.ExpenseFilter
		if err := c.Bind(&params); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"error": err.Error(),
			})
		}

		userID := c.Get("user_id").(string)
		expenses, err := h.service.ListExpense(c.Request().Context(), userID, params)
		if err != nil {
			code, resp := httpres.MapServiceErrHTTPResponse(err)
			return c.JSON(code, resp)
		}

		return c.JSON(http.StatusCreated, map[string]any{
			"data": expenses,
		})
	}
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

func (h *handler) ListExpenseCategory() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := c.Get("user_id").(string)
		categories, err := h.service.ListExpenseCategory(c.Request().Context(), userID)
		if err != nil {
			code, resp := httpres.MapServiceErrHTTPResponse(err)
			return c.JSON(code, resp)
		}

		return c.JSON(http.StatusOK, map[string]any{
			"data": categories,
		})
	}
}
