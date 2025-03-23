package main

import (
	"gokafka/internal/services/notificator/app"
	"os"
)

func main() {
	var env string

	env = os.Getenv("env")
	if env == "" {
		env = "local"
	}

	app, err := app.New(env)
	if err != nil {
		panic(err)
	}

	app.Run()
}
