package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"demo/infrastructure"
)

type ChatController interface {
	ChatSockets(w http.ResponseWriter, r *http.Request)
}

type ChatControllerImpl struct {
	socket *websocket.Upgrader
}

func (c *ChatControllerImpl) ChatSockets(w http.ResponseWriter, r *http.Request) {
	conn, err := c.socket.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("Connection socket error: ", err)
		return
	}
	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Fatal("Read message error: ", err)
			return
		}
		fmt.Printf("================================/n --> message: %d - %s/n", messageType, string(message))
	}
}

// Declare a variable to socket connection
func NewChatController() ChatController {
	var socket = infrastructure.GetSocket()
	return &ChatControllerImpl{
		socket: socket,
	}
}
