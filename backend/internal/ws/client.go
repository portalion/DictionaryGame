package ws

import (
	"time"

	"github.com/gorilla/websocket"
)

var (
	pongWait = 10 * time.Second
	pingInterval = (pongWait * 9) / 10
)

type Client struct {
	connection *websocket.Conn

	writeMessageChannel chan Event
}

func NewClient(connection *websocket.Conn, disconnectChannel chan *Client, incomingMessageHandler chan Event) *Client {
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