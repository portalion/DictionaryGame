package ws

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type RoomManager struct {
	Disconnect       chan *Client
	Connect          chan *websocket.Conn
	IncommingMessage chan Event

	Clients map[*Client] bool

	sync.RWMutex
}

func NewRoomManager() *RoomManager {
	result := &RoomManager{
		Disconnect: make(chan *Client),
		Connect: make(chan *websocket.Conn),
		IncommingMessage: make(chan Event, 8),

		Clients: make(map[*Client]bool),
	}

	go result.run()
	
	return result
}

func (rm *RoomManager) run() {
	for {
		select {
		case wsConnection := <-rm.Connect:
			log.Println("Wanted connection")
			rm.addClient(wsConnection)
		case client := <-rm.Disconnect:
			log.Println("Wanted disconnection")
			rm.removeClient(client)
		case message := <-rm.IncommingMessage:
			log.Println(message.Type)
		}
	}
}

func (rm *RoomManager) removeClient(client *Client) {
	rm.Lock()
	defer rm.Unlock()

	if _, ok := rm.Clients[client]; ok {
		client.connection.Close()
		delete(rm.Clients, client)

		log.Println("Client disconnected")
	}
}

func (rm *RoomManager) addClient(connection *websocket.Conn) {
	rm.Lock()
	defer rm.Unlock()

	client := NewClient(connection, rm.Disconnect, rm.IncommingMessage)
	rm.Clients[client] = true
}