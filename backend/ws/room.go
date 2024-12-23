package ws

import (
	"encoding/json"
	"errors"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

var (


	ErrEventNotSupported = errors.New("this event type is not supported")
)

type Room struct {
	users UserList

	sync.RWMutex
	handlers map[string]EventHandler
}

func NewRoom() *Room {
	m := &Room{
		users: make(UserList),
		handlers: make(map[string]EventHandler),
	}

	m.setupEventHandlers();

	return m;
}

func (m *Room) setupEventHandlers() {
	// m.handlers[UserJoinedMessage] = func(e Event, c *User) error {
	// 	fmt.Println(e)
	// 	return nil
	// }
}

func (m *Room) routeEvent(event Event, c *User) error {
	if handler, ok := m.handlers[event.Type]; ok {
		if err := handler(event, c); err != nil {
			return err
		}
		return nil
	} else {
		return ErrEventNotSupported
	}
}

func (m *Room) JoinRoom(username string, conn *websocket.Conn) {
	user := NewUser(conn, username)
	m.addUser(user)

	go user.readMessages(m)
	go user.writeMessages(m)
}

func (room *Room) addUser(user *User) {
	room.Lock()
	defer room.Unlock()
	
	data, err := json.Marshal(user.username)
	if err != nil {
		log.Println(err)
		user.connection.Close()
		return
	}

	for u := range room.users {
		u.messageChannel <- Event{Type: UserJoinedMessage, Payload: data}
	}
	log.Println("user connected")
	room.users[user] = true
}

func (room *Room) removeUser(user *User) {
	room.Lock()
	defer room.Unlock()

	if _, ok := room.users[user]; ok {
		user.connection.Close()
		delete(room.users, user)

		data, err := json.Marshal(user.username)
		if err != nil {
		}

		log.Println("user disconnected")
		for u := range room.users {
			u.messageChannel <- Event{Type: UserDisconnectedMessage, Payload: data}
	}
	}
}