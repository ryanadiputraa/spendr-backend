package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ryanadiputraa/spendr-backend/config"
	"github.com/ryanadiputraa/spendr-backend/pkg/logger"
)

type Server struct {
	config *config.Config
	log    logger.Logger
	web    *echo.Echo
}

func NewHTTPServer(config *config.Config, log logger.Logger) *Server {
	return &Server{
		config: config,
		log:    log,
		web:    echo.New(),
	}
}

func (s *Server) ServeHTTP() error {
	s.web.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization", "Access-Control-Allow-Origin"},
		AllowMethods: []string{"OPTIONS", "GET", "POST", "PUT", "PATCH", "DELETE"},
	}))
	return s.web.Start(s.config.Server.Port)
}
