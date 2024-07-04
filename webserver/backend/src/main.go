package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
)

const WEBAPP_PATH = "/app/dist/frontend"
const GAME_SERVER_URL = "localhost"
const (
	USER_CHANNEL        = "10000"
	ADMIN_CHANNEL       = "10001"
	SUPER_ADMIN_CHANNEL = "10002"
)

var currId = -1

func main() {
	// fs := http.FileServer(http.Dir(WEBAPP_PATH))

	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	if r.URL.Path != "/" {
	// 		fullPath := WEBAPP_PATH + strings.TrimPrefix(path.Clean(r.URL.Path), "/")
	// 		_, err := os.Stat(fullPath)
	// 		if err != nil {
	// 			if !os.IsNotExist(err) {
	// 				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// 				return
	// 			}
	// 			r.URL.Path = "/"
	// 		}
	// 	}
	// 	fs.ServeHTTP(w, r)
	// })
	// http.HandleFunc("/api", getHello)
	// if err := http.ListenAndServe(getPort(), nil); err != nil {
	// 	log.Fatal(err)
	// }

	//
	// GORILLA MUX
	//

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
	router.HandleFunc("/api/scoreboard", scoreboard).Methods(http.MethodGet)
	router.HandleFunc("/api/game", getParameters).Methods(http.MethodGet)
	router.HandleFunc("/api/game", setParameters).Methods(http.MethodPost)
	router.HandleFunc("/api/status", getStatus).Methods(http.MethodGet)
	router.HandleFunc("/api/start", startGame).Methods(http.MethodPost)
	router.HandleFunc("/api/stop", stopGame).Methods(http.MethodPost)

	// websocket
	router.HandleFunc("/echo", echoWs)
	router.HandleFunc("/ws/game", getGameInfo)
	router.HandleFunc("/ws/scoreboard", getScoreboardInfo)

	// redirects all unhandled paths to the frontend
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("/app/dist/frontend/")))

	port := getPort()
	fmt.Printf("Starting server on port %s\n", port)

	log.Fatal(http.ListenAndServe(port, router))

	//
	// GIN
	//

	// router := gin.Default()

	// router.GET("/api", helloHandler)
	// // redirects all unhandled paths to the frontend
	// router.Static("/static", "/app/dist/frontend")
	// router.NoRoute(func(c *gin.Context) {
	// 	c.File("/app/dist/frontend/index.html")
	// })

	// port := getPort()
	// fmt.Printf("Starting server on port %s\n", port)
	// err := router.Run(port)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Fatal(http.ListenAndServe(port, router))
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
