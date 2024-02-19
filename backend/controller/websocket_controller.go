package controller

import (
	"net/http"

	"github.com/gorilla/websocket"

	"demo/service"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

//================================================================
// Declare a variable to socket connection

type SocketController interface {
	UpgradeChats(w http.ResponseWriter, r *http.Request)
}

type SocketControllerImpl struct {
	chatService service.ChatService
}

func (c *SocketControllerImpl) UpgradeChats(w http.ResponseWriter, r *http.Request) {

}

// Declare a variable to socket connection
func NewChatController() SocketController {
	chatService := service.NewChatService()
	return &SocketControllerImpl{
		chatService: chatService,
	}
}
