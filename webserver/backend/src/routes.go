package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"strconv"
)

type Observer interface {
	update(string)
	getID() string
}

type WebSocket struct {
	id     string
	WSconn *websocket.Conn
}

func (w *WebSocket) update(data string) {
	w.WSconn.WriteMessage(websocket.TextMessage, []byte(data))
}
func (w *WebSocket) getID() string {
	return w.id
}

func getHello(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /api request\n")
	io.WriteString(w, "Hello, HTTP!\n")
}

var upgrader = websocket.Upgrader{}

func echoWs(w http.ResponseWriter, r *http.Request) {
	WSconn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer WSconn.Close()

	for {
		messageType, p, err := WSconn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		if err := WSconn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}

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

func getGameInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got ws://game\n")
	WSconn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer WSconn.Close()

	WSconn.WriteMessage(websocket.TextMessage, []byte(gameBroadcast.data))
	gameBroadcast.register(&WebSocket{WSconn: WSconn, id: "game"})
	defer gameBroadcast.deregister(&WebSocket{WSconn: WSconn, id: "game"})
	for {
		_, _, err := WSconn.ReadMessage()
		if err != nil {
			gameBroadcast.deregister(&WebSocket{WSconn: WSconn, id: "game"})
			return
		}

	}
}

func getScoreboardInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got ws://scoreboard\n")
	WSconn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer WSconn.Close()

	WSconn.WriteMessage(websocket.TextMessage, []byte(scoreBroadcast.data))
	scoreBroadcast.register(&WebSocket{WSconn: WSconn, id: "scoreboard"})
	defer scoreBroadcast.deregister(&WebSocket{WSconn: WSconn, id: "scoreboard"})
	for {
		_, _, err := WSconn.ReadMessage()
		if err != nil {
			scoreBroadcast.deregister(&WebSocket{WSconn: WSconn, id: "scoreboard"})
			return
		}
	}
}
