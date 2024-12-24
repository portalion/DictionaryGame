package room

import (
	. "server/internal/ws/client"
	. "server/internal/ws/event"
	"server/internal/ws/user"

	"log"

	"github.com/gorilla/websocket"
)

type WantsConnectionData struct {
	Username string
	WsConnection *websocket.Conn 
}

type RoomManager struct {
	Disconnect       chan *Client
	Connect          chan WantsConnectionData
	IncommingMessage chan RoomManagerClientMessage

	Clients map[*Client] *user.User
}

func NewRoomManager() *RoomManager {
	result := &RoomManager{
		Disconnect: make(chan *Client),
		Connect: make(chan WantsConnectionData),
		IncommingMessage: make(chan RoomManagerClientMessage, 8),

		Clients: make(map[*Client] *user.User),
	}

	go result.run()
	
	return result
}

func (rm *RoomManager) run() {
	for {
		select {
		case wsConnectionData := <-rm.Connect:
			rm.addClient(wsConnectionData.WsConnection, wsConnectionData.Username)
		case client := <-rm.Disconnect:
			rm.removeClient(client)
		case message := <-rm.IncommingMessage:
			switch message.Request.Type {
			case RoomStateRequested:
				responsePayload := RoomStateEvent{Users: make([]string, len(rm.Clients))}
				i := 0
				for _, user := range rm.Clients {
					responsePayload.Users[i] = user.Username ; i++
				}

				event, err := CreateEvent(RoomState, responsePayload)
				if err != nil {
					log.Println(err)
					continue
				}

				message.Sender.SendEvent(event)
			}
		}
	}
}

func (rm *RoomManager) broadcast(event Event) {
	for client := range rm.Clients{
		client.SendEvent(event)
	}
}

func (rm *RoomManager) removeClient(client *Client) {
	if userData, ok := rm.Clients[client]; ok {
		client.CloseConnection()
		delete(rm.Clients, client)

		event, err := CreateEvent(UserJoinedMessage, 
			UserRelatedMessage{Username: userData.Username,})
		if err != nil {
			log.Println(err)
		}

		rm.broadcast(event)
		log.Println("Client disconnected")
	}
}

func (rm *RoomManager) addClient(connection *websocket.Conn, username string) {
	client := NewClient(connection, rm.Disconnect, rm.IncommingMessage)

	event, err := CreateEvent(UserJoinedMessage, 
		UserRelatedMessage{Username: username,})
	if err != nil {
		log.Println(err)
	}

	rm.broadcast(event)
	rm.Clients[client] = user.NewUser(username, client)

	log.Println("Client connected")
}