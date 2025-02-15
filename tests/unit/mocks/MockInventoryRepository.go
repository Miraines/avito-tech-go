package mocks

import (
	"avito-tech-go/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockInventoryRepository struct {
	mock.Mock
}

func (m *MockInventoryRepository) CreateItem(item *domain.InventoryItem) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MockInventoryRepository) UpdateItem(item *domain.InventoryItem) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MockInventoryRepository) GetByUserAndType(userID uint, itemType string) (*domain.InventoryItem, error) {
	args := m.Called(userID, itemType)
	inv, _ := args.Get(0).(*domain.InventoryItem)
	return inv, args.Error(1)
}

func (m *MockInventoryRepository) GetAllByUser(userID uint) ([]domain.InventoryItem, error) {
	args := m.Called(userID)
	invList, _ := args.Get(0).([]domain.InventoryItem)
	return invList, args.Error(1)
}

func (m *MockInventoryRepository) DeleteItem(item *domain.InventoryItem) error {
	args := m.Called(item)
	return args.Error(0)
}
