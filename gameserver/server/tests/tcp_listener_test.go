package tests

import (
	"net"
	"server/communications"
	"server/interfaces"
	"strconv"
	"testing"
)

func tcpListenerTestRunListener(l interfaces.Listener, hasListenerFinished *bool, t *testing.T, end *chan bool) {
	var err error
	err = l.Run()
	if err != nil {
		t.Fatal(err)
	}
	*hasListenerFinished = true
	*end <- true
}

func tcpListenerTestRunClient(hasClientFinished *bool, t *testing.T, end *chan bool, message string, magic uint32, port uint32) {
	var conn net.Conn
	var err error
	conn, _ = net.Dial("tcp", "127.0.0.1:"+strconv.Itoa(int(port)))

	var header *communications.Header = communications.NewHeaderFromValue(magic, uint32(len(message)))
	_, err = conn.Write(header.AsByte())
	if err != nil {
		t.Fatal(err)
	}

	_, err = conn.Write([]byte(message))
	if err != nil {
		t.Fatal(err)
	}
	err = conn.Close()
	if err != nil {
		t.Fatal(err)
	}

	*hasClientFinished = true
	*end <- true
}

func TestTcpListenerWithOneChannel(t *testing.T) {
	var message string = "allo"
	var magic uint32 = 0x11223344
	var port uint32 = 10000

	var hasListenerFinished bool = false
	var hasClientFinished bool = false
	var listener interfaces.Listener
	var channel chan string = make(chan string)
	defer close(channel)
	var listeningChannel chan bool = make(chan bool)
	var pLC = &listeningChannel
	defer close(listeningChannel)
	listener = interfaces.NewTcpListener(port, &channel, magic)
	listener.AddListeningChannel(pLC)

	go tcpListenerTestRunListener(listener, &hasListenerFinished, t, pLC)
	_ = <-listeningChannel
	go tcpListenerTestRunClient(&hasClientFinished, t, pLC, message, magic, port)
	_ = <-listeningChannel
	listener.Stop()
	var receivedMessage string
	receivedMessage = <-channel
	_ = <-listeningChannel
	if !hasListenerFinished {
		t.Fatalf("hasListenerFinished should be true")
	}
	if !hasClientFinished {
		t.Fatalf("hasClientFinished should be true")
	}
	if receivedMessage != message {
		t.Fatalf("receivedMessage should be %s, but is %s", message, receivedMessage)
	}
	listener.Stop()
}
