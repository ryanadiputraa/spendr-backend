package main

import (
	"fmt"
	"os"

	"github.com/ryanadiputraa/spendr-backend/config"
	"github.com/ryanadiputraa/spendr-backend/internal/server"
	"github.com/ryanadiputraa/spendr-backend/pkg/logger"
)

func main() {
	mode := os.Args[1]

	log, err := logger.NewLogger()
	if err != nil {
		panic(err)
	}

	config, err := config.LoadConfig("env", fmt.Sprintf("./config/.env.%v", mode))
	if err != nil {
		log.Fatal("load config: ", err)
	}

	server := server.NewHTTPServer(config, log)
	if err := server.ServeHTTP(); err != nil {
		log.Fatal("start server: ", err)
	}
}
