package ws

import (
	"log"

	"github.com/gorilla/websocket"
)

type RoomManagerClientMessage struct {
	request Event
	sender *Client
}

type RoomManager struct {
	Disconnect       chan *Client
	Connect          chan *websocket.Conn
	IncommingMessage chan RoomManagerClientMessage

	Clients map[*Client] bool
}

func NewRoomManager() *RoomManager {
	result := &RoomManager{
		Disconnect: make(chan *Client),
		Connect: make(chan *websocket.Conn),
		IncommingMessage: make(chan RoomManagerClientMessage, 8),

		Clients: make(map[*Client]bool),
	}

	go result.run()
	
	return result
}

func (rm *RoomManager) run() {
	for {
		select {
		case wsConnection := <-rm.Connect:
			rm.addClient(wsConnection)
		case client := <-rm.Disconnect:
			rm.removeClient(client)
		case message := <-rm.IncommingMessage:
			log.Println(message.request.Type)
		}
	}
}

func (rm *RoomManager) broadcast(event Event) {
	for client := range rm.Clients{
		client.SendEvent(event)
	}
}

func (rm *RoomManager) removeClient(client *Client) {
	if _, ok := rm.Clients[client]; ok {
		client.connection.Close()
		delete(rm.Clients, client)

		log.Println("Client disconnected")
	}
}

func (rm *RoomManager) addClient(connection *websocket.Conn) {
	client := NewClient(connection, rm.Disconnect, rm.IncommingMessage)
	rm.Clients[client] = true

	log.Println("Client connected")
}