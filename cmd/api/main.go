package main

import (
	"fmt"
	"os"

	"github.com/ryanadiputraa/spendr-backend/config"
)

func main() {
	mode := os.Args[1]

	_, err := config.LoadConfig("env", fmt.Sprintf("./config/.env.%v", mode))
	if err != nil {
		panic(err)
	}
}
