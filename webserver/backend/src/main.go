package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	communications "webserver/Model"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
)

const WEBAPP_PATH = "/app/dist/frontend"

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
	send_command("status")
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

func send_command(command string) string {
	// send request to TCP port 10000 (localhost)
	conn, err := net.Dial("tcp", "localhost:10001")
	if err != nil {
		log.Fatal(err)
	} else {
		var connErr error
		var readLen int

		var comm = communications.NewCommunication(command, 0x11223344)
		conn.Write((comm.AsByte()))
		var buff = make([]byte, 8)
		readLen, connErr = conn.Read(buff)
		if connErr != nil {
			log.Fatal(connErr)
		} else if readLen != 8 {
			log.Fatal("header should be of length 8")
		} else {
			var header, _ = communications.NewHeaderFromBytes(buff)
			var messageBuff = make([]byte, header.GetMessageLength())
			readLen, connErr = conn.Read(messageBuff)
			if connErr != nil {
				log.Fatal(connErr)
			} else {
				return string(messageBuff)
			}
		}

	}
	return ""
}
