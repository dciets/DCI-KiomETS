package communications

import "encoding/json"

type Message struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

func NewMessage(Type string, content string) *Message {
	return &Message{Type: Type, Content: content}
}

func FromJson(jsonStr string) *Message {
	var msg Message
	json.Unmarshal([]byte(jsonStr), &msg)
	return &msg
}

func (m *Message) String() string {
	var msg []byte
	msg, _ = json.Marshal(m)
	return string(msg)
}
