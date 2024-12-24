package event

import (
	"encoding/json"
	"fmt"
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

func CreateEvent(eventType string, payload interface{}) (Event, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return Event{}, fmt.Errorf("failed to marshal payload: %w", err)
	}

	return Event{Type: eventType, Payload: data}, nil
}