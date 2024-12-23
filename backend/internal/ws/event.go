package ws

import "encoding/json"

type Event struct {
	Type string `json:"type"`
	Payload json.RawMessage `json:"payload"`
}
	
const (
	UserJoinedMessage = "user_joined"
	UserDisconnectedMessage = "user_disconnected"
	RoomStateRequest = "room_state_requested"
)