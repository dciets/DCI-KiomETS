package tests

import "errors"

type ListenerMuck struct {
	receiver          *chan string
	closed            bool
	callWhenListening *chan bool
	dataStore         *chan string
}

func (m *ListenerMuck) Run() error {
	if m.closed {
		return errors.New("mock already closed")
	}
	if m.callWhenListening != nil {
		*m.callWhenListening <- true
	}
	v, ok := <-*m.dataStore
	for ok {
		*m.receiver <- v

		v, ok = <-*m.dataStore
	}
	return nil
}

func (m *ListenerMuck) Stop() {
	m.closed = true
	close(*m.dataStore)
}

func (m *ListenerMuck) Write(msg []byte) error {
	if m.closed {
		return errors.New("mock already closed")
	}
	*m.dataStore <- string(msg)
	return nil
}

func (m *ListenerMuck) AddListeningChannel(channel *chan bool) {
	m.callWhenListening = channel
}
