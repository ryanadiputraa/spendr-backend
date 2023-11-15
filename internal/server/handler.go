package server

import (
	_openAPIMiddleware "github.com/go-openapi/runtime/middleware"
	"github.com/labstack/echo/v4"
	"github.com/ryanadiputraa/spendr-backend/internal/auth"
	"github.com/ryanadiputraa/spendr-backend/internal/expense"
	"github.com/ryanadiputraa/spendr-backend/internal/middleware"
	"github.com/ryanadiputraa/spendr-backend/internal/user"
	"github.com/ryanadiputraa/spendr-backend/pkg/jwt"
	"github.com/ryanadiputraa/spendr-backend/pkg/validator"
)

func (s *Server) setupHandlers() {
	authGroup := s.web.Group("/auth")
	userGroup := s.web.Group("/api/users")
	expenseGroup := s.web.Group("/api/expenses")

	validator := validator.NewValidator()
	jwtService := jwt.NewJWTService()
	authMiddleware := middleware.NewAuthMiddleware(s.config, jwtService)

	userRepository := user.NewRepository(s.db)
	userService := user.NewService(s.log, userRepository)
	user.NewHandler(userGroup, userService, *authMiddleware)

	authService := auth.NewService(s.log, validator, userRepository)
	auth.NewHandler(authGroup, s.config, s.log, validator, authService, jwtService)

	expenseRepository := expense.NewRepository(s.db)
	expenseService := expense.NewService(s.log, expenseRepository)
	expense.NewHandler(expenseGroup, validator, expenseService, *authMiddleware)

	sh := _openAPIMiddleware.Redoc(_openAPIMiddleware.RedocOpts{
		SpecURL: "/api/open-api/open-api.yml",
		Path:    "/api/docs",
		Title:   "Spendr - API Docs",
	}, nil)

	s.web.Static("/api/open-api", "api/open-api")
	s.web.GET("/api/docs", func(c echo.Context) error {
		sh.ServeHTTP(c.Response().Writer, c.Request())
		return nil
	})
}
