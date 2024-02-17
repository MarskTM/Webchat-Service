package infrastructure

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	socket = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins
		},
	}
)


func GetSocket() *websocket.Upgrader {
	return &socket
}

func init() {}