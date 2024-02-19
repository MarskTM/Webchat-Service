package repository

import (
	"demo/infrastructure"
	"demo/model"

	"gorm.io/gorm"
)

type HubRepository interface {
	GetAll() ([]*model.HubTable, error)
	CreateHub(record *model.HubTable) (*model.HubTable, error)
	UpdateHub(record *model.HubTable) (*model.HubTable, error)
	DeleteHub(record *model.HubTable) (*model.HubTable, error)
}

type HubRepositoryImpl struct {
	db *gorm.DB
}

func (h *HubRepositoryImpl) GetAll() ([]*model.HubTable, error) {
	var records []*model.HubTable
	if err := h.db.Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

func (h *HubRepositoryImpl) CreateHub(record *model.HubTable) (*model.HubTable, error) {
	if err := h.db.Create(record).Error; err != nil {
		return nil, err
	}
	return record, nil
}

func (h *HubRepositoryImpl) UpdateHub(record *model.HubTable) (*model.HubTable, error) {
	if err := h.db.Save(record).Error; err != nil {
		return nil, err
	}
	return record, nil
}

func (h *HubRepositoryImpl) DeleteHub(record *model.HubTable) (*model.HubTable, error) {
	if err := h.db.Delete(record).Error; err != nil {
		return nil, err
	}
	return record, nil
}

// Sử dụng design pattern Factory để tạo ra một đối tượng HubRepository
// Nhằm mục đích tạo ra một đối tượng HubRepository mà không cần biết cụ thể là đối tượng nào
func NewHubRepository() HubRepository {
	db := infrastructure.GetDB()
	return &HubRepositoryImpl{
		db: db,
	}
}
