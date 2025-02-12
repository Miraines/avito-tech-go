package repositories

import (
	"avito-tech-go/internal/domain"
	"errors"
	"gorm.io/gorm"
)

type InventoryRepository interface {
	CreateItem(item *domain.InventoryItem) error
	UpdateItem(item *domain.InventoryItem) error
	GetByUserAndType(userID uint, itemType string) (*domain.InventoryItem, error)
	GetAllByUser(userID uint) ([]domain.InventoryItem, error)
	DeleteItem(item *domain.InventoryItem) error
}

type inventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) InventoryRepository {
	return &inventoryRepository{db: db}
}

func (r *inventoryRepository) CreateItem(item *domain.InventoryItem) error {
	return r.db.Create(item).Error
}

func (r *inventoryRepository) UpdateItem(item *domain.InventoryItem) error {
	return r.db.Save(item).Error
}

func (r *inventoryRepository) GetByUserAndType(userID uint, itemType string) (*domain.InventoryItem, error) {
	var inv domain.InventoryItem
	err := r.db.Where("user_id = ? AND item_type = ?", userID, itemType).First(&inv).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &inv, err
}

func (r *inventoryRepository) GetAllByUser(userID uint) ([]domain.InventoryItem, error) {
	var items []domain.InventoryItem
	err := r.db.Where("user_id = ?", userID).Find(&items).Error
	return items, err
}

func (r *inventoryRepository) DeleteItem(item *domain.InventoryItem) error {
	return r.db.Delete(item).Error
}
