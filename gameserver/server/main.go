package main

import (
	"github.com/joho/godotenv"
	"os"
	"server/app"
)

func loadEnv() {
	if os.Getenv("APP_ENV") == "" {
		var err error = godotenv.Load()
		if err != nil {
			panic("Error loading .env file")
		}
	}
}

func main() {
	loadEnv()

	var application *app.Application
	var appErr error
	application, appErr = app.LoadApplicationFromEnv()
	if appErr != nil {
		panic(appErr)
	}

}
