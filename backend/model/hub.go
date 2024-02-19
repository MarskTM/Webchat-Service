package model

import "gorm.io/gorm"

type Hub struct {
	Name       string           // Tên của hu
	clients    map[*Client]bool // Danh sách các máy khách đang hoạt động
	broadcast  chan []byte      // Tin nhắn gửi đến từ máy khách.
	register   chan *Client     // Đăng ký máy khách mới.
	unregister chan *Client     // Hủy đăng ký máy khách.
}

// Quản lý số lượng hub trong hệ thống
type HubTable struct {
	gorm.Model // Thêm các trường ID, CreatedAt, UpdatedAt, DeletedAt

	Name   string `json:"name"`   // Tên của hub
	Status bool   `json:"status"` // Trạng thái của hub
}



