package state

import (
	"log"
	"server/internal/ws/event"
	"server/internal/ws/user"
	"slices"
)

type RoomState struct {
	Users []*user.User

	broadcast func(event.Event)
	swapState func(State)
}

func NewRoomState(users []*user.User, broadcast func(event.Event), swapState func(State)) *RoomState{
	return &RoomState{
		Users: users,
		broadcast: broadcast,
		swapState: swapState,
	}
}

func (rs *RoomState) OnUserConnection(user *user.User) {
	rs.Users = append(rs.Users, user)
}

func (rs *RoomState) OnUserDisconnection(user *user.User) {
	i := slices.Index(rs.Users, user)
	rs.Users = slices.Delete(rs.Users, i, i)
}

func (rs *RoomState) handleStateRequestedEvent(sender *user.User) error {
	responsePayload := event.RoomStateEvent{Users: make([]string, len(rs.Users))}
	i := 0
	for _, user := range rs.Users {
		responsePayload.Users[i] = user.Username ; i++
	}
	event, err := event.CreateEvent(event.RoomState, responsePayload)
	if err != nil {
		return err
	}
	sender.ConnectionClient.SendEvent(event)
	return nil
}

func (rs *RoomState) ProcessMessage(message ServerStateMessage) error {
	switch message.Event.Type {
	case event.RoomStateRequested:
		return rs.handleStateRequestedEvent(message.Sender)
	}

	log.Printf("Got message from %s\n", message.Sender.Username)
	return nil
}