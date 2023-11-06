package main

import (
	"fmt"
	"os"

	"github.com/ryanadiputraa/spendr-backend/config"
	"github.com/ryanadiputraa/spendr-backend/pkg/logger"
)

func main() {
	mode := os.Args[1]

	log, err := logger.NewLogger()
	if err != nil {
		panic(err)
	}

	_, err = config.LoadConfig("env", fmt.Sprintf("./config/.env.%v", mode))
	if err != nil {
		log.Fatal(err)
	}

	log.Info("starting")
}
