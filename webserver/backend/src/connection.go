package main

import (
	"log"
	"net"
	"os"
	"sync"
	"time"
	"webserver/Model/communications"
)

type MutexChanQueue struct {
	mu    sync.Mutex
	queue []chan string
}

func (q *MutexChanQueue) Push(c chan string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.queue = append(q.queue, c)
}

func (q *MutexChanQueue) Pop() chan string {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.queue) == 0 {
		return nil
	}
	var c = q.queue[0]
	q.queue = q.queue[1:]
	return c
}

type clientConnObj struct {
	conn       net.Conn
	game       chan string
	scoreboard chan string
}
type adminQueueObj struct {
	conn     net.Conn
	channels MutexChanQueue
}
type Connection struct {
	clientConn clientConnObj
	adminQueue adminQueueObj
}

var lock = &sync.Mutex{}

// create a singleton connection
var conn *Connection

func GetConnection() *Connection {
	var gamePortStr string = os.Getenv("GAME_PORT")
	var adminPortStr string = os.Getenv("ADMIN_PORT")
	var gameServerHost string = os.Getenv("GAME_SERVER_HOST")
	if conn == nil {
		lock.Lock()
		defer lock.Unlock()
		if conn == nil {
			conn = &Connection{}
			var err error
			conn.clientConn.conn, err = net.Dial("tcp", gameServerHost+":"+gamePortStr)
			if err != nil {
				log.Fatal(err)
			}
			conn.adminQueue.conn, err = net.Dial("tcp", gameServerHost+":"+adminPortStr)
			if err != nil {
				log.Fatal(err)
			}
			conn.clientConn.game = make(chan string)
			conn.clientConn.scoreboard = make(chan string)
			// start the listener
			go conn.ReadClient()
			go conn.SendCommand()
		}
	}
	return conn
}

func (c *Connection) SendCommand() {
	var connErr error
	var readLen int
	for {
		if len(c.adminQueue.channels.queue) != 0 {
			channel := c.adminQueue.channels.Pop()
			command := <-channel
			buffer := communications.NewCommunication(command, 0x11223344).AsByte()
			log.Printf("sending command : %s", command)
			_, err := c.adminQueue.conn.Write(buffer)
			if err != nil {
				log.Fatal(err)
			}
			var buff = make([]byte, 8)
			readLen, connErr = c.adminQueue.conn.Read(buff)
			if connErr != nil {
				log.Fatal(connErr)
			}
			if readLen != 8 {
				log.Fatal("header should be of length 8")
			}
			var header, _ = communications.NewHeaderFromBytes(buff)
			var messageBuff = make([]byte, header.GetMessageLength())
			if header.GetMessageLength() > 0 {
				readLen, connErr = c.adminQueue.conn.Read(messageBuff)
				if connErr != nil {
					log.Fatal(connErr)
				}
			}
			log.Printf("got result : `%s`", string(messageBuff))
			channel <- string(messageBuff)
		} else {
			time.Sleep(50 * time.Millisecond)
		}
	}
}

func fullRead(header *communications.Header, conn net.Conn) ([]byte, error) {
	var length uint32 = header.GetMessageLength()
	var reed uint32 = 0
	var bytes []byte = make([]byte, length)

	for reed < length {
		var currentRead []byte = make([]byte, length-reed)
		var reedLength int
		var err error
		reedLength, err = conn.Read(currentRead)
		if err != nil {
			return nil, err
		}

		copy(bytes[reed:reed+uint32(reedLength)], currentRead[0:reedLength])
		reed += uint32(reedLength)
	}

	return bytes, nil
}

func (c *Connection) ReadClient() {
	//c.clientConn.game <- "{'message': 'Game hasn't started yet'}"
	//c.clientConn.scoreboard <- "{'message': 'Game hasn't started yet'}"
	for {
		var connErr error
		var readLen int
		var buff = make([]byte, 8)
		readLen, connErr = c.clientConn.conn.Read(buff)
		if connErr != nil {
			log.Fatal(connErr)
		}
		if readLen != 8 {
			log.Fatal("header should be of length 8")
		}
		var header, _ = communications.NewHeaderFromBytes(buff)

		var messageBuff []byte
		messageBuff, connErr = fullRead(header, c.clientConn.conn)
		if connErr != nil {
			log.Fatal(connErr)
		}

		msg := communications.FromJson(string(messageBuff))
		if msg.Type == "scoreboard" {
			c.clientConn.scoreboard <- msg.Content
		} else {
			c.clientConn.game <- msg.Content
		}
	}
}
