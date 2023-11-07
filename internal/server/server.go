package server

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ryanadiputraa/spendr-backend/config"
	"github.com/ryanadiputraa/spendr-backend/internal/auth"
	"github.com/ryanadiputraa/spendr-backend/internal/user"
	"github.com/ryanadiputraa/spendr-backend/pkg/jwt"
	"github.com/ryanadiputraa/spendr-backend/pkg/logger"
	"github.com/ryanadiputraa/spendr-backend/pkg/validator"
)

type Server struct {
	config *config.Config
	log    logger.Logger
	web    *echo.Echo
	db     *sqlx.DB
}

func NewHTTPServer(config *config.Config, log logger.Logger, db *sqlx.DB) *Server {
	return &Server{
		config: config,
		log:    log,
		web:    echo.New(),
		db:     db,
	}
}

func (s *Server) ServeHTTP() error {
	s.setupHandlers()
	s.web.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization", "Access-Control-Allow-Origin"},
		AllowMethods: []string{"OPTIONS", "GET", "POST", "PUT", "PATCH", "DELETE"},
	}))
	return s.web.Start(s.config.Server.Port)
}

func (s *Server) setupHandlers() {
	authGroup := s.web.Group("/auth")

	validator := validator.NewValidator()
	jwtService := jwt.NewJWTService()

	userRepository := user.NewRepository(s.db)

	authService := auth.NewService(s.log, validator, userRepository)
	auth.NewHandler(authGroup, s.config, s.log, validator, authService, jwtService)
}
