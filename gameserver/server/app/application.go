package app

import (
	"errors"
	"os"
	"strconv"
)

type Application struct {
	AdminPort uint32
	GamePort  uint32
}

func NewApplication(adminPort uint32, gamePort uint32) *Application {
	return &Application{AdminPort: adminPort, GamePort: gamePort}
}

func LoadApplicationFromEnv() (*Application, error) {
	var gamePortStr string = os.Getenv("GAME_PORT")
	if gamePortStr == "" {
		return nil, errors.New("GAME_PORT not set")
	}

	var gamePort int
	var gamePortErr error
	gamePort, gamePortErr = strconv.Atoi(gamePortStr)
	if gamePortErr != nil {
		return nil, errors.New("GAME_PORT is not a number")
	}

	var adminPortStr string = os.Getenv("ADMIN_PORT")
	if adminPortStr == "" {
		return nil, errors.New("ADMIN_PORT not set")
	}

	var adminPort int
	var adminPortErr error
	adminPort, adminPortErr = strconv.Atoi(adminPortStr)
	if adminPortErr != nil {
		return nil, errors.New("ADMIN_PORT is not a number")
	}

	var app *Application = NewApplication(uint32(adminPort), uint32(gamePort))

	return app, nil
}
