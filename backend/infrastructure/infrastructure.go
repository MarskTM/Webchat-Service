package infrastructure

import (
	"sync"

	"github.com/gorilla/websocket"
)

var (
	// Initialize a mutex to safely access the connections slice.
	connectionsMutex sync.Mutex

	// Initialize an empty slice to hold active connections.
	connections []*websocket.Conn
)

func GetSocketConnections() []*websocket.Conn {
	return connections
}

func GetSocketConnectionsMutex() *sync.Mutex {
	return &connectionsMutex
}

func init() {}
