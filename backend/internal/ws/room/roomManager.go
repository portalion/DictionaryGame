package room

import (
	. "server/internal/ws/client"
	. "server/internal/ws/event"
	"server/internal/ws/state"
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
	State state.State
}

func NewRoomManager() *RoomManager {
	result := &RoomManager{
		Disconnect: make(chan *Client),
		Connect: make(chan WantsConnectionData),
		IncommingMessage: make(chan RoomManagerClientMessage, 8),

		Clients: make(map[*Client] *user.User),
		State: &state.RoomState{},
	}

	go result.run()
	
	return result
}

func (rm *RoomManager) ChangeState(state state.State) {
	rm.State = state
}

func (rm *RoomManager) run() {
	for {
		select {
		case wsConnectionData := <-rm.Connect:
			user := rm.addClient(wsConnectionData.WsConnection, wsConnectionData.Username)
			rm.State.OnUserConnection(user)
		case client := <-rm.Disconnect:
			user := rm.removeClient(client)
			if user != nil {
				rm.State.OnUserDisconnection(user)
			}
		case message := <-rm.IncommingMessage:
			err := rm.State.ProcessMessage(
				state.ServerStateMessage{
					Sender: rm.Clients[message.Sender],
					Event: message.Request},
				rm.broadcast,
				rm.ChangeState)

			if err != nil {
				log.Println(err)
				continue
			}
		}
	}
}

func (rm *RoomManager) broadcast(event Event) {
	for client := range rm.Clients{
		client.SendEvent(event)
	}
}

func (rm *RoomManager) removeClient(client *Client) *user.User{
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
		return userData
	}
	return nil
}

func (rm *RoomManager) addClient(connection *websocket.Conn, username string) *user.User{
	client := NewClient(connection, rm.Disconnect, rm.IncommingMessage)

	event, err := CreateEvent(UserJoinedMessage, 
		UserRelatedMessage{Username: username,})
	if err != nil {
		log.Println(err)
	}

	rm.broadcast(event)
	rm.Clients[client] = user.NewUser(username, client)

	log.Println("Client connected")
	return rm.Clients[client]
}

