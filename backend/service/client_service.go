package service

import (
	"bytes"
	"demo/model"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second    // Thời gian chờ phản hồi từ client
	pingPeriod     = (pongWait * 9) / 10 // Thời gian gửi ping
	maxMessageSize = 512                 // Kích thước tin nhăn tối đa là 512 bytes
)

// Sử dụng biến này để thay thế kí tự xuống dòng
var (
	newline = []byte("\n")
	space   = []byte(" ")
)

// Client sẽ là đối tượng trung gian giữa websocket connection và hub
type Client struct {
	hub    *Hub            // Dia chỉ Hub sẽ quản lý các client
	conn   *websocket.Conn // Đối tượng websocket connection
	send   chan []byte     // Kênh dùng để gửi tin nhắn từ hub tới websocket connection
	roomId int             // Tên phòng chat
}

// ----------------------------------------------------------------------------

type ChatServiceImpl interface{}

// Đọc tin nhắn từ websocket connection tới hub
func (s *Client) readPump() {

	defer func() {
		s.hub.unregister <- c // Gửi thông báo hủy đăng ký client từ hub
		s.conn.Close()        // Đóng kết nối
	}()

	s.conn.SetReadLimit(maxMessageSize)              //	Đặt giới hạn kích thước tin nhắn
	s.conn.SetReadDeadline(time.Now().Add(pongWait)) // Đặt thời gian chờ phản hồi từ client
	s.conn.SetPongHandler(func(string) error {
		s.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	}) // Đặt hàm xử lý khi nhận được tin nhắn pong từ client

	// Dọc tin nhắn từ websocket connection
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			// websocket.IsUnexpectedCloseError dùng để kiểm tra xem lỗi có phải là lỗi đóng kết nối không
			// websocket.CloseGoingAway dùng để thông báo rằng client đã đóng kết nối
			// websocket.CloseAbnormalClosure dùng để thông báo rằng kết nối bị đóng không đúng cách
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				// Ghi log lỗi
				log.Println("Error:", err)
			}
			break
		}

		// TrimSpace dùng để loại bỏ khoảng trắng ở đầu và cuối chuỗi
		// -1 dùng để loại bỏ tất cả khoảng trắng thừa
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		s.hub.broadcast <- message // Gửi tin nhắn tới hub
	}
}

// Gửi tin nhắn từ hub tới websocket connection
func (s *Client) writePump() {

	// Ticker là một đối tượng dùng để gửi tin nhắn ping tới client
	/*
		Việc sử dụng ticker nhằm mục dích kiểm tra xem client có đang hoạt động hay không, nếu không hoạt động thì sẽ đóng kết nối.
	*/
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		s.conn.Close()
	}()

	// Cú pháp for - select dùng để lắng nghe từ nhiều kênh khác nhau
	for {
		select {
		// Thực hiện gửi tin nhắn nếu trong kênh c.send có tin nhắn
		case message, ok := <-s.send:
			s.conn.SetWriteDeadline(time.Now().Add(writeWait)) // Đặt thời gian chờ ghi tin nhắn

			// Nếu kênh c.send bị đóng thì thoát khỏi vòng lặp
			/*
				Sử dụng conn.WriteMessage để gửi tin nhắn thay vì c.conn.NextWriter lý do là
				c.conn.NextWriter sẽ gửi tin nhắn theo kiểu streaming trong khi conn.WriteMessage sẽ gửi tin nhắn theo kiểu block:
					- Nhược điểm của kiểu block là nó sẽ chờ đến khi gửi tin nhắn xong mới tiếp tục thực hiện các lệnh tiếp theo.
					- Còn kiểu streaming sẽ gửi tin nhắn ngay lập tức và tiếp tục thực hiện các lệnh tiếp theo.
			*/
			if !ok {
				s.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// Sửu dụng đối tượng w dùng để ghi tin nhắn.
			w, err := s.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message) // Ghi tin nhắn vào đối tượng w

			// Gửi toàn bộ tin nhắn từ kênh c.send tới websocket connection
			n := len(s.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-s.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		// Thực hiện gửi tin nhắn ping tới client dựa trênh channel ticker.C (được khởi tạo từ đối tượng ticker)
		case <-ticker.C:
			s.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := s.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func NewChatService() ChatService {
	return &ChatServiceImpl{}
}
