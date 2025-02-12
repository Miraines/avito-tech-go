package repositories

import (
	"avito-tech-go/internal/domain"
	"errors"

	"gorm.io/gorm"
)

type MerchRepository interface {
	CreateMerchItem(item *domain.MerchItem) error
	UpdateMerchItem(item *domain.MerchItem) error
	GetMerchItemByType(itemType string) (*domain.MerchItem, error)
	GetAllMerchItems() ([]domain.MerchItem, error)
	DeleteMerchItem(itemType string) error
}

type merchRepository struct {
	db *gorm.DB
}

func NewMerchRepository(db *gorm.DB) MerchRepository {
	return &merchRepository{db: db}
}

func (r *merchRepository) CreateMerchItem(item *domain.MerchItem) error {
	return r.db.Create(item).Error
}

func (r *merchRepository) UpdateMerchItem(item *domain.MerchItem) error {
	return r.db.Save(item).Error
}

func (r *merchRepository) GetMerchItemByType(itemType string) (*domain.MerchItem, error) {
	var item domain.MerchItem
	err := r.db.Where("item_type = ?", itemType).First(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &item, err
}

func (r *merchRepository) GetAllMerchItems() ([]domain.MerchItem, error) {
	var items []domain.MerchItem
	err := r.db.Find(&items).Error
	return items, err
}

func (r *merchRepository) DeleteMerchItem(itemType string) error {
	res := r.db.Where("item_type = ?", itemType).Delete(&domain.MerchItem{})
	return res.Error
}
