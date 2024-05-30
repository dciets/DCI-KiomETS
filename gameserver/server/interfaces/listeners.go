package interfaces

type Listener interface {
	Run() error
	Stop()
	Write([]byte) error
	AddListeningChannel(*chan bool)
}
