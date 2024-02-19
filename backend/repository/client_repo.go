package repository

import (
	"demo/infrastructure"
	"demo/model"

	"gorm.io/gorm"
)

type ClientRepository struct {
	db *gorm.DB
}

type ClientRepositoryImpl interface {
	CreateClient(record *model.CLientLog) (*model.CLientLog, error)
	UpdateClient(record *model.CLientLog) (*model.CLientLog, error)
	DeleteClient(record *model.CLientLog) (*model.CLientLog, error)
	GetAllClient() ([]*model.CLientLog, error)
}

func (c *ClientRepository) GetAllClient() ([]*model.CLientLog, error) {
	var records []*model.CLientLog
	if err := c.db.Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

func (c *ClientRepository) CreateClient(record *model.CLientLog) (*model.CLientLog, error) {
	if err := c.db.Create(record).Error; err != nil {
		return nil, err
	}
	return record, nil
}

func (c *ClientRepository) UpdateClient(record *model.CLientLog) (*model.CLientLog, error) {
	if err := c.db.Save(record).Error; err!= nil {
        return nil, err
    }
    return record, nil
}

func (c *ClientRepository) DeleteClient(record *model.CLientLog) (*model.CLientLog, error) {
	if err := c.db.Delete(record).Error; err!= nil {
        return nil, err
    }
    return record, nil
}

// Sử dụng design pattern Factory để tạo ra một đối tượng ClientRepository
// Nhằm mục đích tạo ra một đối tượng ClientRepository mà không cần biết cụ thể là đối tượng nào
func NewClientRepository() ClientRepositoryImpl {
	db := infrastructure.GetDB()
	return &ClientRepository{
		db: db,
	}
}
