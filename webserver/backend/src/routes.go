package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"webserver/Model"
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

func init() {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
}

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
	channel <- "all-player " + strconv.Itoa(currId) + "-admin_command"
	result := <-channel
	//var encoded = strings.Split(result, " ")[1]
	var encoded = strings.Split(result, " ")[1]
	var decoded, err = base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	//create agents from decoded and maping id to UID
	type MidleAgent struct {
		UID  string `json:"id"`
		Name string `json:"name"`
	}
	var midleAgent []MidleAgent
	var err2 = json.Unmarshal(decoded, &midleAgent)
	if err2 != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	var agents = make([]Model.Agent, len(midleAgent))
	for i := range agents {
		agents[i].UID = midleAgent[i].UID
		agents[i].Name = midleAgent[i].Name
	}
	// encode in json and send
	var toSend, err3 = json.Marshal(agents)
	if err3 != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(toSend)
}

func createAgent(w http.ResponseWriter, r *http.Request) {
	var conn = GetConnection()
	channel := make(chan string)
	conn.adminQueue.channels.Push(channel)
	// read the boduy
	var agent Model.Agent
	// Decode the incoming JSON request body and check for errors
	err := json.NewDecoder(r.Body).Decode(&agent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Now you can access agent.Name
	agentName := agent.Name
	log.Printf("Received agent name: %s\n", agentName)
	// get the name of the agent and encode it in b64
	var encodedName = base64.StdEncoding.EncodeToString([]byte(agent.Name))
	channel <- "new-player " + strconv.Itoa(currId) + "-" + RandStringRunes(10) + " " + encodedName
	result := <-channel
	var data = strings.Split(result, " ")
	if len(data) < 2 {
		http.Error(w, "Agent name already exists", http.StatusBadRequest)
		return
	}
	// return the UID of the agent
	agent.UID = data[1]
	// decode the UID in b64
	var decodedUID, err2 = base64.StdEncoding.DecodeString(agent.UID)
	if err2 != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	agent.UID = string(decodedUID)
	// encode in json and send
	toSend, err3 := json.Marshal(agent)
	if err3 != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(toSend)
}

func getParameters(w http.ResponseWriter, r *http.Request) {
	var conn = GetConnection()
	channel := make(chan string)
	conn.adminQueue.channels.Push(channel)
	channel <- "get-parameters " + strconv.Itoa(currId) + "-admin_command"
	result := <-channel
	// decode the parameters
	var encoded = strings.Split(result, " ")[1]
	var decoded, err = base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(decoded)
}
func setParameters(w http.ResponseWriter, r *http.Request) {
	// read the request body for the parameters
	var data, err = io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	// encode the parameters in b64
	var encoded = base64.StdEncoding.EncodeToString(data)
	var conn = GetConnection()
	channel := make(chan string)
	conn.adminQueue.channels.Push(channel)
	channel <- "set-parameters " + strconv.Itoa(currId) + "-admin_command " + encoded
	_ = <-channel
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func startGame(w http.ResponseWriter, r *http.Request) {
	var conn = GetConnection()
	channel := make(chan string)
	conn.adminQueue.channels.Push(channel)
	channel <- "start " + strconv.Itoa(currId) + "-admin_command"
	result := <-channel
	// get int answer and convert to int
	var data = strings.Split(result, " ")
	var returnVal, err = strconv.Atoi(data[1])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if returnVal != 1 {
		http.Error(w, "Game already started", http.StatusConflict)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("Game started"))
}
func stopGame(w http.ResponseWriter, r *http.Request) {
	var conn = GetConnection()
	channel := make(chan string)
	conn.adminQueue.channels.Push(channel)
	channel <- "stop " + strconv.Itoa(currId) + "-admin_command"
	result := <-channel
	// get int answer and convert to int
	var data = strings.Split(result, " ")
	var returnVal, err = strconv.Atoi(data[1])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if returnVal != 1 {
		http.Error(w, "Game isn't started", http.StatusConflict)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("Game will stop after the current round ends"))
}

func getStatus(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /api/status request\n")
	var conn = GetConnection()
	channel := make(chan string)
	conn.adminQueue.channels.Push(channel)
	channel <- "status " + strconv.Itoa(currId) + "-admin_command"
	result := <-channel
	// get int answer and convert to int
	var data = strings.Split(result, " ")
	var returnVal, err = strconv.Atoi(data[1])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	if returnVal != 1 {
		w.Write([]byte("Game is not running"))
	} else {
		w.Write([]byte("Game is currently running"))
	}
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
