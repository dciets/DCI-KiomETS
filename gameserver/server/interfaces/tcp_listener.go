package interfaces

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"server/communications"
	"time"
)

type TcpListener struct {
	port              uint32
	receiver          *chan string
	isTcpClosed       bool
	tcpListener       net.Listener
	connections       []net.Conn
	timeBeforeTimeout time.Time
	callWhenListening *chan bool
	magic             uint32
}

func NewTcpListener(port uint32, receiver *chan string, timeBeforeTimeout time.Time, magic uint32) *TcpListener {
	return &TcpListener{port: port, receiver: receiver, timeBeforeTimeout: timeBeforeTimeout, isTcpClosed: true, connections: []net.Conn{}, tcpListener: nil, callWhenListening: nil, magic: magic}
}

func (t *TcpListener) AddListeningChannel(channel *chan bool) {
	t.callWhenListening = channel
}

func (t *TcpListener) Stop() {
	t.isTcpClosed = true
	_ = t.tcpListener.Close()
}

func (t *TcpListener) Run() error {
	if !t.isTcpClosed {
		return errors.New("TcpListener is already running")
	}
	var tcpListener net.Listener
	var tcpError error
	tcpListener, tcpError = net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", t.port))
	if tcpError != nil {
		return tcpError
	}
	t.isTcpClosed = false
	t.tcpListener = tcpListener
	defer t.Stop()
	if t.callWhenListening != nil {
		*t.callWhenListening <- true
	}

	for !t.isTcpClosed {
		var conn net.Conn
		conn, tcpError = t.tcpListener.Accept()
		if tcpError == nil {
			t.connections = append(t.connections, conn)
			_ = conn.SetReadDeadline(t.timeBeforeTimeout)
			go t.acceptRequest(conn)
		}
	}

	return nil
}

func (t *TcpListener) Write(message []byte) error {
	if t.isTcpClosed {
		return errors.New("listener is closed")
	}
	var hasError bool = false
	for _, conn := range t.connections {
		_, err := conn.Write(message)
		if err != nil {
			hasError = true
		}
	}

	if hasError {
		return errors.New("write error")
	}
	return nil
}

func (t *TcpListener) acceptRequest(conn net.Conn) {
	defer conn.Close()

	for !t.isTcpClosed {
		var connErr error
		var protocol []byte = make([]byte, 8)
		var readByte int

		// Lire plus de byte
		readByte, connErr = conn.Read(protocol)
		if connErr != nil {
			log.Fatal(connErr)
		}
		if errors.Is(connErr, io.ErrClosedPipe) {
			return
		}
		if readByte < 8 {
			continue
		}

		var header *communications.Header
		header, _ = communications.NewHeaderFromBytes(protocol)

		if header.GetHeaderMagic() != t.magic {
			continue
		}

		var messageBuff []byte = make([]byte, header.GetMessageLength())
		readByte, connErr = conn.Read(messageBuff)

		if uint32(readByte) != header.GetMessageLength() {
			continue
		}

		var message = string(messageBuff[:])
		*t.receiver <- message
	}
}
