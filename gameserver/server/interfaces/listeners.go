package interfaces

type Listener interface {
	Run() error
	Stop()
	Write(string) error
	AddListeningChannel(*chan bool)
}
