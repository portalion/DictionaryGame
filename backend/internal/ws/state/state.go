package state

import (
	"server/internal/ws/event"
	"server/internal/ws/user"
)

type ServerStateMessage struct {
	Sender *user.User
	Event event.Event
}

type State interface {
	OnUserConnection(*user.User)
	OnUserDisconnection(*user.User)
	ProcessMessage(message ServerStateMessage) error
}