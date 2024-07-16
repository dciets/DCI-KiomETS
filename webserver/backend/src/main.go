package main

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
)

var currId = -1

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
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With ", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{os.Getenv("ORIGIN_ALLOWED")})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	var conn = GetConnection()
	channel := make(chan string)
	conn.adminQueue.channels.Push(channel)
	channel <- "id-assign"
	result := <-channel

	currId, _ = strconv.Atoi(result[3:])

	startChannelReading()

	router := mux.NewRouter()

	router.HandleFunc("/api", getHello).Methods(http.MethodGet)

	router.HandleFunc("/api/agent", getAgents).Methods(http.MethodGet)
	router.HandleFunc("/api/agent", createAgent).Methods(http.MethodPost)

	router.HandleFunc("/api/game", getParameters).Methods(http.MethodGet)
	router.HandleFunc("/api/game", setParameters).Methods(http.MethodPut)

	router.HandleFunc("/api/start", startGame).Methods(http.MethodPost)
	router.HandleFunc("/api/stop", stopGame).Methods(http.MethodPost)

	router.HandleFunc("/api/status", getStatus).Methods(http.MethodGet)

	// websocket
	router.HandleFunc("/echo", echoWs)
	router.HandleFunc("/ws/game", getGameInfo)
	router.HandleFunc("/ws/scoreboard", getScoreboardInfo)

	port := getPort()
	fmt.Printf("Starting server on port %s\n", port)

	log.Fatal(http.ListenAndServe(port, handlers.CORS(originsOk, headersOk, methodsOk)(router)))
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// Returns the port from the "PORT" environment variable, or returns ":8080" by default.
func getPort() string {
	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = ":8080"
	}
	return port
}

func startChannelReading() {
	var conn = GetConnection()
	go func() {
		for {
			gameData := <-conn.clientConn.game
			gameBroadcast.setData(gameData)
		}
	}()
	go func() {
		for {
			scoreData := <-conn.clientConn.scoreboard
			scoreBroadcast.setData(scoreData)
		}
	}()
}
