package app

import (
	"errors"
	"log"
	"os"
	"server/interfaces"
	"strconv"
)

type Application struct {
	SuperAdminPort   uint32
	AdminPort        uint32
	GamePort         uint32
	UserRepository   *UserRepository
	GameRuntime      *GameRuntime
	CommandProcessor *CommandProcessor
	CommChannel      *chan string
	DuplicateChannel *chan string
}

func NewApplication(adminPort uint32, gamePort uint32, superAdminPort uint32) *Application {
	var commChannel chan string = make(chan string)
	var duplicateChannel = make(chan string)
	return &Application{
		AdminPort:        adminPort,
		GamePort:         gamePort,
		SuperAdminPort:   superAdminPort,
		CommChannel:      &commChannel,
		DuplicateChannel: &duplicateChannel,
		UserRepository:   nil,
		GameRuntime:      nil,
		CommandProcessor: nil,
	}
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

	var superAdminPortStr string = os.Getenv("SUPER_ADMIN_PORT")
	if superAdminPortStr == "" {
		return nil, errors.New("SUPER_ADMIN_PORT not set")
	}

	var superAdminPort int
	var superAdminPortErr error
	superAdminPort, superAdminPortErr = strconv.Atoi(superAdminPortStr)
	if superAdminPortErr != nil {
		return nil, errors.New("SUPER_ADMIN_PORT is not a number")
	}

	var app *Application = NewApplication(uint32(adminPort), uint32(gamePort), uint32(superAdminPort))

	return app, nil
}

func (application *Application) Start(applicationBootstrapChannel *chan bool) {
	defer close(*application.CommChannel)
	defer close(*application.DuplicateChannel)

	var superAdminListener interfaces.Listener = interfaces.NewTcpListener(application.SuperAdminPort, application.CommChannel, 0x11223344)
	var adminListener interfaces.Listener = interfaces.NewTcpListener(application.AdminPort, application.CommChannel, 0x11223344)
	var clientListener interfaces.Listener = interfaces.NewTcpListener(application.GamePort, application.CommChannel, 0x11223344)

	var superAdminChannelListener interfaces.Listener = interfaces.NewChannelListener(superAdminListener, application.DuplicateChannel)
	var adminChannelListener interfaces.Listener = interfaces.NewChannelListener(adminListener, application.DuplicateChannel)
	var clientChannelListener interfaces.Listener = interfaces.NewChannelListener(clientListener, application.DuplicateChannel)

	var superAdminListenerRunning = true
	var adminListenerRunning = true
	var clientListenerRunning = true

	go func() {
		err := superAdminChannelListener.Run()
		if err != nil {
			log.Print(err)
		}
		superAdminListenerRunning = false
	}()
	go func() {
		err := adminChannelListener.Run()
		if err != nil {
			log.Print(err)
		}
		adminListenerRunning = false
	}()
	go func() {
		err := clientChannelListener.Run()
		if err != nil {
			log.Print(err)
		}
		clientListenerRunning = false
	}()

	application.UserRepository = NewUserRepository()
	application.GameRuntime = NewGameRuntime(clientChannelListener)
	application.CommandProcessor = NewCommandProcessor(application.UserRepository, application.GameRuntime, clientChannelListener, adminChannelListener, superAdminChannelListener)

	if applicationBootstrapChannel != nil {
		*applicationBootstrapChannel <- true
	}

	for superAdminListenerRunning && adminListenerRunning && clientListenerRunning {
		var message string
		message = <-*application.CommChannel
		application.CommandProcessor.Process(message)
	}
}
