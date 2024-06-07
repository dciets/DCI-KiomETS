package interfaces

type ChannelListener struct {
	otherListener Listener
	duplicate     *chan string
}

func NewChannelListener(otherListener Listener, duplicate *chan string) *ChannelListener {
	return &ChannelListener{
		otherListener: otherListener,
		duplicate:     duplicate,
	}
}

func (c *ChannelListener) Run() error {
	return c.otherListener.Run()
}

func (c *ChannelListener) Stop() {
	c.otherListener.Stop()
}

func (c *ChannelListener) Write(s string) error {
	*c.duplicate <- s
	return c.otherListener.Write(s)
}

func (c *ChannelListener) AddListeningChannel(*chan bool) {

}
