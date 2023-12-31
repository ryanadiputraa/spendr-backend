package main

import (
	config "github.com/ryanadiputraa/spendr-backend/configs"
	"github.com/ryanadiputraa/spendr-backend/internal/server"
	"github.com/ryanadiputraa/spendr-backend/pkg/db/postgres"
	"github.com/ryanadiputraa/spendr-backend/pkg/logger"
)

func main() {
	log, err := logger.NewLogger()
	if err != nil {
		panic(err)
	}

	config, err := config.LoadConfig("yml", "./configs/config.yml")
	if err != nil {
		log.Fatal("load config: ", err)
	}

	db, err := postgres.NewDB(config)
	if err != nil {
		log.Fatal("open postgres conn: ", err)
	}

	server := server.NewHTTPServer(config, log, db)
	if err := server.ServeHTTP(); err != nil {
		log.Fatal("start server: ", err)
	}
}
