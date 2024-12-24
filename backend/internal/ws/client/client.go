package client

import (
	. "server/internal/ws/event"
	"time"

	"github.com/gorilla/websocket"
)

var (
	pongWait = 10 * time.Second
	pingInterval = (pongWait * 9) / 10
)

type RoomManagerClientMessage struct {
	Request Event
	Sender *Client
}

type Client struct {
	connection *websocket.Conn
	writeMessageChannel chan Event
}

func NewClient(connection *websocket.Conn, disconnectChannel chan *Client, incomingMessageHandler chan RoomManagerClientMessage) *Client {
	result := &Client{
		connection: connection,
		writeMessageChannel: make(chan Event),
	}

	go result.handleWriting(disconnectChannel)
	go result.handleReading(disconnectChannel, incomingMessageHandler)

	return result
}

func (client *Client) SendEvent(event Event) {
	client.writeMessageChannel <- event
}

func (client *Client) CloseConnection() {
	client.connection.Close()
}