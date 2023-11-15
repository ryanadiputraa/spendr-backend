package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	_echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/ryanadiputraa/spendr-backend/config"

	"github.com/ryanadiputraa/spendr-backend/pkg/logger"
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
	s.web.Use(_echoMiddleware.CORSWithConfig(_echoMiddleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization", "Access-Control-Allow-Origin"},
		AllowMethods: []string{"OPTIONS", "GET", "POST", "PUT", "PATCH", "DELETE"},
	}))

	server := &http.Server{
		Addr:         s.config.Server.Port,
		Handler:      s.web,
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			s.log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, os.Kill)

	sig := <-quit
	s.log.Info("received terminate, graceful shutdown ", sig)

	tc, shutdown := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdown()

	return server.Shutdown(tc)
}
