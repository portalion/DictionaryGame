package client

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"

	"server/internal/ws/event"
)

func (client *Client) dispatchMessage(message event.Event, incomingMessageHandler chan RoomManagerClientMessage) {
	incomingMessageHandler <- RoomManagerClientMessage{Request: message, Sender: client}
}

func (client *Client) pongHandler(pongMessage string) error {
	return client.connection.SetReadDeadline(time.Now().Add(pongWait))
}

func (client *Client) handleReading(disconnect chan *Client, incomingMessageHandler chan RoomManagerClientMessage) {
	defer func() {
		disconnect <- client
	}()

	if err := client.connection.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Println(err)
		return
	}

	client.connection.SetPongHandler(client.pongHandler)

	for {
		_, payload, err := client.connection.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNoStatusReceived, websocket.CloseNormalClosure) {
				log.Printf("error reading message: %v", err)
			}
			break
		}

		var request event.Event
		if err := json.Unmarshal(payload, &request); err != nil {
			log.Printf("error marshalling message: %v", err)
			break
		}

		client.dispatchMessage(request, incomingMessageHandler)
	}
}

func (client *Client) handleWriting(disconnect chan *Client) {
	ticker := time.NewTicker(pingInterval)
	defer func() {
		ticker.Stop()
		disconnect <- client
	}()

	for {
		select {
		case message, ok := <-client.writeMessageChannel:
			if !ok {
				if err := client.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("connection closed: ", err)
				}
				return
			}

			data, err := json.Marshal(message)
			if err != nil {
				log.Println(err)
				return
			}

			if err := client.connection.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Println(err)
			}
		case <-ticker.C:
			if err := client.connection.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}
