package main

import (
	"log"
	"net"
	"sync"
	"webserver/Model/communications"
)

type connObj struct {
	conn    net.Conn
	channel chan string
}

type Connection struct {
	admin  connObj
	client connObj
}

var lock = &sync.Mutex{}

// create a singleton connection
var conn *Connection

func GetConnection() *Connection {
	if conn == nil {
		lock.Lock()
		defer lock.Unlock()
		if conn == nil {
			conn = &Connection{}
			var err error
			conn.admin.conn, err = net.Dial("tcp", "localhost:10000")
			if err != nil {
				log.Fatal(err)
			}
			conn.admin.channel = make(chan string)
			conn.client.conn, err = net.Dial("tcp", "localhost:10001")
			if err != nil {
				log.Fatal(err)
			}
			conn.client.channel = make(chan string)

			// start the listener
			go conn.ReadChannel("admin")
			go conn.ReadChannel("client")
		}
	}
	return conn
}

func (c *Connection) SendCommand(command string) string {
	var connErr error
	var readLen int
	var comm = communications.NewCommunication(command, 0x11223344)
	c.admin.conn.Write((comm.AsByte()))
	var buff = make([]byte, 8)
	readLen, connErr = c.admin.conn.Read(buff)
	if connErr != nil {
		log.Fatal(connErr)
	} else if readLen != 8 {
		log.Fatal("header should be of length 8")
	}
	var header, _ = communications.NewHeaderFromBytes(buff)
	var messageBuff = make([]byte, header.GetMessageLength())
	readLen, connErr = c.admin.conn.Read(messageBuff)
	if connErr != nil {
		log.Fatal(connErr)
	}

	return string(messageBuff)
}

func (c *Connection) ReadChannel(channel string) {
	var connection connObj
	switch channel {
	case "admin":
		connection = c.admin
	case "client":
		connection = c.client
	}
	var connErr error
	var readLen int
	var buff = make([]byte, 8)
	readLen, connErr = connection.conn.Read(buff)
	if connErr != nil {
		log.Fatal(connErr)
	}
	if readLen != 8 {
		log.Fatal("header should be of length 8")
	}
	var header, _ = communications.NewHeaderFromBytes(buff)
	var messageBuff = make([]byte, header.GetMessageLength())
	readLen, connErr = connection.conn.Read(messageBuff)
	if connErr != nil {
		log.Fatal(connErr)
	}
	connection.channel <- string(messageBuff)
}
