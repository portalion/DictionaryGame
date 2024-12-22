package ws

import "encoding/json"

type Event struct {
	Type string `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type EventHandler func(event Event, c *User) error
	
const (
	UserJoinedMessage = "user_joined"
	UserDisconnectedMessage = "user_disconnected"
)