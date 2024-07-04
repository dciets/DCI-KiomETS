package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
	"io"
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

func getAgents(w http.ResponseWriter, r *http.Request) {
	var conn = GetConnection()
	channel := make(chan string)
	conn.adminQueue.channels.Push(channel)
	channel <- "all-players " + strconv.Itoa(currId) + "-admin_command"
	result := <-channel
	fmt.Fprintf(w, result)
}

func createAgent(w http.ResponseWriter, r *http.Request) {
	var conn = GetConnection()
	channel := make(chan string)
	conn.adminQueue.channels.Push(channel)
	channel <- "new-player " + strconv.Itoa(currId) + "-" + RandStringRunes(10)
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

func getParameters(w http.ResponseWriter, r *http.Request) {
	var conn = GetConnection()
	channel := make(chan string)
	conn.adminQueue.channels.Push(channel)
	channel <- "get-parameters " + strconv.Itoa(currId) + "-admin_command"
	result := <-channel
	fmt.Fprintf(w, result)
}
func setParameters(w http.ResponseWriter, r *http.Request) {
	// read the request body for the parameters
	var data, err = io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("got data: %s\n", string(data))
	var conn = GetConnection()
	channel := make(chan string)
	conn.adminQueue.channels.Push(channel)
	channel <- "set-parameters " + strconv.Itoa(currId) + "-admin_command " + string(data)
	result := <-channel
	fmt.Fprintf(w, result)
}

func startGame(w http.ResponseWriter, r *http.Request) {
	var conn = GetConnection()
	channel := make(chan string)
	conn.adminQueue.channels.Push(channel)
	channel <- "start " + strconv.Itoa(currId) + "-admin_command"
	result := <-channel
	fmt.Fprintf(w, result)
}
func stopGame(w http.ResponseWriter, r *http.Request) {
	var conn = GetConnection()
	channel := make(chan string)
	conn.adminQueue.channels.Push(channel)
	channel <- "stop " + strconv.Itoa(currId) + "-admin_command"
	result := <-channel
	fmt.Fprintf(w, result)
}

func getStatus(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /api/status request\n")
	var conn = GetConnection()
	channel := make(chan string)
	conn.adminQueue.channels.Push(channel)
	channel <- "status " + strconv.Itoa(currId) + "-admin_command"
	result := <-channel
	fmt.Fprintf(w, result)
}
