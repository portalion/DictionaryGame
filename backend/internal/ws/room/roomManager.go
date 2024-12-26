package room

import (
	"server/internal/ws/client"
	"server/internal/ws/event"
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
	Disconnect       chan *client.Client
	Connect          chan WantsConnectionData
	IncommingMessage chan client.RoomManagerClientMessage

	Clients map[*client.Client] *user.User
	State state.State
}

func NewRoomManager() *RoomManager {
	result := &RoomManager{
		Disconnect: make(chan *client.Client),
		Connect: make(chan WantsConnectionData),
		IncommingMessage: make(chan client.RoomManagerClientMessage, 8),

		Clients: make(map[*client.Client] *user.User),
	}

	result.State = state.NewRoomState(make([]*user.User, 0), result.broadcast, result.ChangeState)
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
					Event: message.Request})

			if err != nil {
				log.Println(err)
				continue
			}
		}
	}
}

func (rm *RoomManager) broadcast(event event.Event) {
	for client := range rm.Clients{
		client.SendEvent(event)
	}
}

func (rm *RoomManager) removeClient(client *client.Client) *user.User{
	if userData, ok := rm.Clients[client]; ok {
		client.CloseConnection()
		delete(rm.Clients, client)

		event, err := event.CreateEvent(event.UserDisconnectedMessage, 
			event.UserRelatedMessage{Username: userData.Username,})
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
	client := client.NewClient(connection, rm.Disconnect, rm.IncommingMessage)

	event, err := event.CreateEvent(event.UserJoinedMessage, 
		event.UserRelatedMessage{Username: username,})
	if err != nil {
		log.Println(err)
	}

	rm.broadcast(event)
	rm.Clients[client] = user.NewUser(username, client)

	log.Println("Client connected")
	return rm.Clients[client]
}

