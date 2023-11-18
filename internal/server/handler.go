package server

import (
	_openAPIMiddleware "github.com/go-openapi/runtime/middleware"
	"github.com/labstack/echo/v4"
	_authHandler "github.com/ryanadiputraa/spendr-backend/internal/auth/handler"
	_authService "github.com/ryanadiputraa/spendr-backend/internal/auth/service"
	_expenseHandler "github.com/ryanadiputraa/spendr-backend/internal/expense/Handler"
	_expenseRepository "github.com/ryanadiputraa/spendr-backend/internal/expense/repository"
	_expenseService "github.com/ryanadiputraa/spendr-backend/internal/expense/service"
	"github.com/ryanadiputraa/spendr-backend/internal/middleware"
	_userHandler "github.com/ryanadiputraa/spendr-backend/internal/user/handler"
	_userRepository "github.com/ryanadiputraa/spendr-backend/internal/user/repository"
	_userService "github.com/ryanadiputraa/spendr-backend/internal/user/service"
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

	userRepository := _userRepository.NewRepository(s.db)
	userService := _userService.NewService(s.log, userRepository)
	_userHandler.NewHandler(userGroup, userService, *authMiddleware)

	authService := _authService.NewService(s.log, validator, userRepository)
	_authHandler.NewHandler(authGroup, s.config, s.log, validator, authService, jwtService)

	expenseRepository := _expenseRepository.NewRepository(s.db)
	expenseService := _expenseService.NewService(s.log, expenseRepository)
	_expenseHandler.NewHandler(expenseGroup, validator, expenseService, *authMiddleware)

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
