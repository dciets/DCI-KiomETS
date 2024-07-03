package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
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

	router := mux.NewRouter()

	router.HandleFunc("/api", getHello).Methods(http.MethodGet)
	router.HandleFunc("/api/bot", getBots).Methods(http.MethodGet)
	router.HandleFunc("/api/bot", createBot).Methods(http.MethodPost)
	router.HandleFunc("/api/scoreboard", scoreboard).Methods(http.MethodGet)
	router.HandleFunc("/api/game", getGame).Methods(http.MethodGet)
	router.HandleFunc("/api/game", setGame).Methods(http.MethodPost)
	router.HandleFunc("/api/status", getStatus).Methods(http.MethodGet)
	router.HandleFunc("/api/start", startGame).Methods(http.MethodPost)
	router.HandleFunc("/api/stop", stopGame).Methods(http.MethodPost)

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

func getHello(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /api request\n")
	io.WriteString(w, "Hello, HTTP!\n")
}

func helloHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, HTTP!",
	})
}

// Returns the port from the "PORT" environment variable, or returns ":8080" by default.
func getPort() string {
	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = ":8080"
	}
	return port
}

func getBots(w http.ResponseWriter, r *http.Request) {
	var conn = GetConnection()
	channel := make(chan string)
	conn.adminQueue.channels.Push(channel)
	channel <- "get-players"
	result := <-channel
	fmt.Fprintf(w, result)
}

func createBot(w http.ResponseWriter, r *http.Request) {
	var conn = GetConnection()
	channel := make(chan string)
	conn.adminQueue.channels.Push(channel)
	channel <- "create-player"
	result := <-channel
	fmt.Fprintf(w, result)
}

func scoreboard(w http.ResponseWriter, r *http.Request) {
	var conn = GetConnection()
	channel := make(chan string)
	conn.adminQueue.channels.Push(channel)
	channel <- "get-scoreboard"
	result := <-channel
	fmt.Fprintf(w, result)
}

func getGame(w http.ResponseWriter, r *http.Request) {
	var conn = GetConnection()
	channel := make(chan string)
	conn.adminQueue.channels.Push(channel)
	channel <- "get-parameters"
	result := <-channel
	fmt.Fprintf(w, result)
}
func setGame(w http.ResponseWriter, r *http.Request) {
	var conn = GetConnection()
	channel := make(chan string)
	conn.adminQueue.channels.Push(channel)
	channel <- "set-parameters"
	result := <-channel
	fmt.Fprintf(w, result)
}

func startGame(w http.ResponseWriter, r *http.Request) {
	var conn = GetConnection()
	channel := make(chan string)
	conn.adminQueue.channels.Push(channel)
	channel <- "start"
	result := <-channel
	fmt.Fprintf(w, result)
}
func stopGame(w http.ResponseWriter, r *http.Request) {
	var conn = GetConnection()
	channel := make(chan string)
	conn.adminQueue.channels.Push(channel)
	channel <- "stop"
	result := <-channel
	fmt.Fprintf(w, result)
}

func getStatus(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /api/status request\n")
	var conn = GetConnection()
	channel := make(chan string)
	conn.adminQueue.channels.Push(channel)
	channel <- "status 1 abc"
	result := <-channel
	fmt.Fprintf(w, result)
}
