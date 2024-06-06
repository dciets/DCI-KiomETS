package app

import (
	"errors"
	"os"
	"server/interfaces"
	"strconv"
	"time"
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

func (app *Application) Start() {
	var adminChannel chan string = make(chan string)
	var gameChannel chan string = make(chan string)
	var adminListener interfaces.Listener = interfaces.NewTcpListener(app.AdminPort, &adminChannel, time.Duration(0), 0x11111111)
	var gameListener interfaces.Listener = interfaces.NewTcpListener(app.GamePort, &gameChannel, time.Duration(0), 0x22222222)
	defer adminListener.Stop()
	defer gameListener.Stop()
	defer close(adminChannel)
	defer close(gameChannel)

}
