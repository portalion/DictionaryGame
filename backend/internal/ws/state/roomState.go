package state

import (
	"log"
	"server/internal/ws/event"
	"server/internal/ws/user"
	"slices"
)

type RoomState struct {
	Users []*user.User
}

func (rs *RoomState) OnUserConnection(user *user.User) {
	rs.Users = append(rs.Users, user)
}

func (rs *RoomState) OnUserDisconnection(user *user.User) {
	i := slices.Index(rs.Users, user)
	rs.Users = slices.Delete(rs.Users, i, i)
}

func (rs *RoomState) ProcessMessage(message ServerStateMessage, broadcast func (event.Event), swapState func(State)) error {
	switch message.Event.Type {
	case event.RoomStateRequested:
		responsePayload := event.RoomStateEvent{Users: make([]string, len(rs.Users))}
		i := 0
		for _, user := range rs.Users {
			responsePayload.Users[i] = user.Username ; i++
		}
		event, err := event.CreateEvent(event.RoomState, responsePayload)
		if err != nil {
			return err
		}
		message.Sender.ConnectionClient.SendEvent(event)
	}

	log.Printf("Got message from %s\n", message.Sender.Username)
	return nil
}