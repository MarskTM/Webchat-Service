package model

// Hub duy trì tập hợp các máy khách đang hoạt động và phát các tin nhắn đến tất cả máy khách
type Hub struct {
	// Danh sách các máy khách đang hoạt động
	clients map[*Client]bool

	// Tin nhắn gửi đến từ máy khách.
	broadcast chan []byte

	// Đăng ký máy khách mới.
	register chan *Client

	// Hủy đăng ký máy khách.
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Khởi tạo hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}

		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
