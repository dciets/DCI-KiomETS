package communications

import (
	"encoding/json"
	"log"
)

type Message struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

func NewMessage(Type string, content string) *Message {
	return &Message{Type: Type, Content: content}
}

func FromJson(jsonStr string) *Message {
	var msg Message
	var err error = json.Unmarshal([]byte(jsonStr), &msg)
	if err != nil {
		log.Fatal(err)
	}
	return &msg
}

func (m *Message) String() string {
	var msg []byte
	msg, _ = json.Marshal(m)
	return string(msg)
}
