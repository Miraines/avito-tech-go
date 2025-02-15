package mocks

import (
	"avito-tech-go/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockMerchRepository struct {
	mock.Mock
}

func (m *MockMerchRepository) CreateMerchItem(item *domain.MerchItem) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MockMerchRepository) UpdateMerchItem(item *domain.MerchItem) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MockMerchRepository) GetMerchItemByType(itemType string) (*domain.MerchItem, error) {
	args := m.Called(itemType)
	merch, _ := args.Get(0).(*domain.MerchItem)
	return merch, args.Error(1)
}

func (m *MockMerchRepository) GetAllMerchItems() ([]domain.MerchItem, error) {
	args := m.Called()
	merchList, _ := args.Get(0).([]domain.MerchItem)
	return merchList, args.Error(1)
}

func (m *MockMerchRepository) DeleteMerchItem(itemType string) error {
	args := m.Called(itemType)
	return args.Error(0)
}
