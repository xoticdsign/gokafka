package main

import (
	"os"

	"gokafka/internal/services/api/app"
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
