package tests

import (
	"server/app"
	"testing"
	"time"
)

func TestEnv(t *testing.T) {
	var err error = loadEnv()
	if err != nil {
		t.Fatalf("Error loading .env file")
	}

	var application *app.Application
	var appErr error
	application, appErr = app.LoadApplicationFromEnv()
	if appErr != nil {
		t.Fatalf("Error loading application from .env file")
	}

	if application.AdminPort != 10100 {
		t.Fatalf("Error loading admin port from .env file")
	}

	if application.GamePort != 10000 {
		t.Fatalf("Error loading game port from .env file")
	}

	if application.SuperAdminPort != 10101 {
		t.Fatalf("Error loading super admin port from .env file")
	}
}

func TestApplicationIsWorking(t *testing.T) {
	_ = loadEnv()
	var application *app.Application
	application, _ = app.LoadApplicationFromEnv()
	var eventChannel chan bool = make(chan bool)
	defer close(eventChannel)
	application.GamePort = 20000
	application.SuperAdminPort = 20202
	application.AdminPort = 20200

	go application.Start(&eventChannel)
	<-eventChannel
	var commChannel *chan string = application.CommChannel
	var duplicateChannel *chan string = application.DuplicateChannel

	var ret string

	*commChannel <- "set-max-tick 2"
	*commChannel <- "set-time-per-tick 1000"

	*commChannel <- "start"

	ret = <-*duplicateChannel
	if ret != "1" {
		t.Log(ret)
		t.Fatal("Start error")
	}

	<-time.NewTimer(time.Duration(10) * time.Millisecond).C

	*commChannel <- "stop"

	ret = <-*duplicateChannel
	if ret != "1" {
		t.Log(ret)
		t.Fatal("Stop error")
	}

	ret = <-*duplicateChannel
	if ret != "[]" {
		t.Log(ret)
		t.Fatal("Update error")
	}

	ret = <-*duplicateChannel
	if ret != "[]" {
		t.Log(ret)
		t.Fatal("Update error")
	}
}
