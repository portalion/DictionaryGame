package ws

import (
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {return true},
	}

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

func (m *Room) JoinRoom(w http.ResponseWriter, r *http.Request) {
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	user := NewUser(conn)
	m.addUser(user)

	go user.readMessages(m)
	go user.writeMessages(m)
}

func (room *Room) addUser(user *User) {
	room.Lock()
	defer room.Unlock()

	for u := range room.users {
		u.messageChannel <- Event{Type: UserJoinedMessage}
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

		log.Println("user disconnected")
		for u := range room.users {
			u.messageChannel <- Event{Type: UserDisconnectedMessage}
	}
	}
}