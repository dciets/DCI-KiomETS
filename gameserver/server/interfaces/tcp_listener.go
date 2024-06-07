package interfaces

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"server/communications"
	"slices"
	"sync"
)

type TcpListener struct {
	port              uint32
	receiver          *chan string
	isTcpClosed       bool
	tcpListener       net.Listener
	connections       []net.Conn
	callWhenListening *chan bool
	magic             uint32
	mut               sync.Mutex
}

func NewTcpListener(port uint32, receiver *chan string, magic uint32) *TcpListener {
	return &TcpListener{port: port, receiver: receiver, isTcpClosed: true, connections: []net.Conn{}, tcpListener: nil, callWhenListening: nil, magic: magic}
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
			go t.acceptRequest(conn)
		}
	}

	return nil
}

func (t *TcpListener) Write(message string) error {
	if t.isTcpClosed {
		return errors.New("listener is closed")
	}
	var hasError bool = false
	var communication *communications.Communication = communications.NewCommunication(message, t.magic)
	for _, conn := range t.connections {
		_, err := conn.Write(communication.AsByte())
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
	defer func(conn net.Conn, t *TcpListener) {
		t.mut.Lock()
		_ = conn.Close()

		t.connections = append(t.connections[:slices.Index(t.connections, conn)], t.connections[slices.Index(t.connections, conn)+1:]...)
		t.mut.Unlock()
	}(conn, t)

	for !t.isTcpClosed {
		var connErr error
		var protocol []byte = make([]byte, 8)
		var readByte int

		t.mut.Lock()
		readByte, connErr = conn.Read(protocol)

		if errors.Is(connErr, io.ErrClosedPipe) {
			t.mut.Unlock()
			return
		}
		if errors.Is(connErr, os.ErrDeadlineExceeded) {
			t.mut.Unlock()
			continue
		} else if connErr != nil {
			t.mut.Unlock()
			log.Fatal(connErr)
		}

		if readByte < 8 {
			t.mut.Unlock()
			continue
		}

		var header *communications.Header
		header, _ = communications.NewHeaderFromBytes(protocol)

		if header.GetHeaderMagic() != t.magic {
			t.mut.Unlock()
			continue
		}

		var messageBuff []byte = make([]byte, header.GetMessageLength())

		readByte, connErr = conn.Read(messageBuff)
		t.mut.Unlock()

		if errors.Is(connErr, io.ErrClosedPipe) {
			return
		}

		if errors.Is(connErr, os.ErrDeadlineExceeded) {
			continue
		} else if connErr != nil {
			log.Fatal(connErr)
		}

		if uint32(readByte) != header.GetMessageLength() {
			continue
		}

		var message = string(messageBuff[:])
		*t.receiver <- message
	}
}
