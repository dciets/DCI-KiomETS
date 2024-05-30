package tests

import (
	"server/app"
	"testing"
)

func TestAllo(t *testing.T) {
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

}
