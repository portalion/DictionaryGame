package event

import (
	"encoding/json"
)

type Event struct {
	Type string `json:"type"`
	Payload json.RawMessage `json:"payload"`
}
	
const (
	UserJoinedMessage = "user_joined"
	UserDisconnectedMessage = "user_disconnected"
	RoomStateRequested = "room_state_requested"
	RoomState = "room_state"
)

type RoomStateEvent struct{
	Users []string `json:"users"` 
}

type UserRelatedMessage struct {
	Username string `json:"username"`
}