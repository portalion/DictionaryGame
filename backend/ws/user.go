package ws

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

var (
	pongWait = 10 * time.Second
	pingInterval = (pongWait * 9) / 10
)

type UserList map[*User]bool

type User struct {
	connection *websocket.Conn

	messageChannel chan Event
	username string
}

func NewUser(connection *websocket.Conn, username string) *User {
	return &User{
		connection: connection,
		messageChannel: make(chan Event),
		username: username,
	}
}

func (user *User) readMessages(roomManager *Room) {
	defer func() {
		roomManager.removeUser(user)
	}()

	if err := user.connection.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Println(err)
		return
	}

	user.connection.SetPongHandler(user.pongHandler)

	for {
		_, payload, err := user.connection.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}
			break
		}

		var request Event
		if err := json.Unmarshal(payload, &request); err != nil {
			log.Printf("error marshalling message: %v", err)
			break 
		}

		if err := roomManager.routeEvent(request, user); err != nil {
			log.Println("Error handeling Message: ", err)
		}
	}
}

func (user *User) writeMessages(roomManager *Room) {
	ticker := time.NewTicker(pingInterval)
	defer func() {
		ticker.Stop()
		roomManager.removeUser(user)
	}()

	for {
		select {
		case message, ok := <-user.messageChannel:
			if !ok {
				if err := user.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("connection closed: ", err)
				}
				return
			}

			data, err := json.Marshal(message)
			if err != nil {
				log.Println(err)
				return
			}

			if err := user.connection.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Println(err)
			}
		case <-ticker.C:
			if err := user.connection.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Println("writemsg: ", err)
				return 
			}
		}

	}
}

func (user *User) pongHandler(pongMsg string) error {
	return user.connection.SetReadDeadline(time.Now().Add(pongWait))
}