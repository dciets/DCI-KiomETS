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

	var updateString string = "{\"type\":\"action\",\"content\":\"{\\\"players\\\":[],\\\"terrains\\\":[{\\\"terrainType\\\":2,\\\"ownerIndex\\\":-1,\\\"numberOfSoldier\\\":0,\\\"position\\\":[0,0]},{\\\"terrainType\\\":2,\\\"ownerIndex\\\":-1,\\\"numberOfSoldier\\\":0,\\\"position\\\":[1,0]},{\\\"terrainType\\\":2,\\\"ownerIndex\\\":-1,\\\"numberOfSoldier\\\":0,\\\"position\\\":[0.5,0.866]},{\\\"terrainType\\\":2,\\\"ownerIndex\\\":-1,\\\"numberOfSoldier\\\":0,\\\"position\\\":[-0.5,0.866]},{\\\"terrainType\\\":2,\\\"ownerIndex\\\":-1,\\\"numberOfSoldier\\\":0,\\\"position\\\":[-1,0]},{\\\"terrainType\\\":2,\\\"ownerIndex\\\":-1,\\\"numberOfSoldier\\\":0,\\\"position\\\":[-0.5,-0.866]},{\\\"terrainType\\\":2,\\\"ownerIndex\\\":-1,\\\"numberOfSoldier\\\":0,\\\"position\\\":[0.5,-0.866]}],\\\"pipes\\\":[{\\\"length\\\":2,\\\"first\\\":0,\\\"second\\\":1,\\\"soldiers\\\":[]},{\\\"length\\\":2,\\\"first\\\":0,\\\"second\\\":2,\\\"soldiers\\\":[]},{\\\"length\\\":2,\\\"first\\\":0,\\\"second\\\":3,\\\"soldiers\\\":[]},{\\\"length\\\":2,\\\"first\\\":0,\\\"second\\\":4,\\\"soldiers\\\":[]},{\\\"length\\\":2,\\\"first\\\":0,\\\"second\\\":5,\\\"soldiers\\\":[]},{\\\"length\\\":2,\\\"first\\\":0,\\\"second\\\":6,\\\"soldiers\\\":[]}]}\"}"
	ret = <-*duplicateChannel
	if ret != updateString {
		t.Log(ret)
		t.Log(updateString)
		t.Fatal("Update error")
	}

	ret = <-*duplicateChannel
	if ret != updateString {
		t.Log(ret)
		t.Log(updateString)
		t.Fatal("Update error")
	}

	ret = <-*duplicateChannel
	if ret != "{\"type\":\"end\",\"content\":\"\"}" {
		t.Log(ret)
		t.Fatal("Update error")
	}
}
