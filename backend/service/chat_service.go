package service

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type ChatService interface {
	PingScoket(conn *websocket.Conn)
}

type ChatServiceImpl struct{}

func (c *ChatServiceImpl) PingScoket(conn *websocket.Conn) {
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			// If an error occurs while reading, assume the connection is closed.
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				log.Println("Connection closed by client:", err)
			} else {
				log.Println("Read message error:", err)
			}
			break
		}
		fmt.Printf("\n--> message: %d - %s", messageType, string(message))

		// Write the message back to the connection.
		msResponse := "Pong"
		conn.WriteMessage(websocket.TextMessage, []byte(msResponse))
	}
}

func NewChatService() ChatService {
	return &ChatServiceImpl{}
}
