package Model

import (
	"fmt"
)

type Agent struct {
	UID   string
	Name  string
	score int
}

func NewAgent(UID string, Name string) *Agent {
	return &Agent{UID: UID, Name: Name}
}

func (b *Agent) String() string {
	return fmt.Sprintf("Agent name : %s, UID : %s", b.Name, b.UID)
}
