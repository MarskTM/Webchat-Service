package controller

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"demo/infrastructure"
	"demo/service"
)

type ChatController interface {
	UpgradeHandler(w http.ResponseWriter, r *http.Request)
}

type ChatControllerImpl struct {
	chatService service.ChatService
}

func (c *ChatControllerImpl) UpgradeHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
		log.Fatal("Connection socket error: ", err)
		return
	}
	defer conn.Close()

	connections := infrastructure.GetSocketConnections()
	connectionsMutex := infrastructure.GetSocketConnectionsMutex()

	// Add the new connection to the list of active connections.
	connectionsMutex.Lock()
	connections = append(connections, conn)
	connectionsMutex.Unlock()

	// Call the service to handle the connection.
	c.chatService.PingScoket(conn)

	// Remove the connection from the list of active connections.
	connectionsMutex.Lock()
	for i, c := range connections {
		if c == conn {
			connections = append(connections[:i], connections[i+1:]...)
			break
		}
	}
	connectionsMutex.Unlock()
}

// Declare a variable to socket connection
func NewChatController() ChatController {
	chatService := service.NewChatService()
	return &ChatControllerImpl{
		chatService: chatService,
	}
}
