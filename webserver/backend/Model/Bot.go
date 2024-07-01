package Model

import (
	"fmt"
)

type Bot struct {
	UID  string
	Name string
}

func NewBot(UID string, Name string) *Bot {
	return &Bot{UID: UID, Name: Name}
}

func (b *Bot) String() string {
	return fmt.Sprintf("Bot name : %s, UID : %s", b.Name, b.UID)
}
