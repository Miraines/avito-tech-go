package mocks

import (
	"avito-tech-go/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) CreateTransaction(tx *domain.Transaction) error {
	args := m.Called(tx)
	return args.Error(0)
}

func (m *MockTransactionRepository) GetUserTransactions(userID uint) ([]domain.Transaction, error) {
	args := m.Called(userID)
	txList, _ := args.Get(0).([]domain.Transaction)
	return txList, args.Error(1)
}

func (m *MockTransactionRepository) GetTransactionsByType(userID uint, txType domain.TransactionType) ([]domain.Transaction, error) {
	args := m.Called(userID, txType)
	txList, _ := args.Get(0).([]domain.Transaction)
	return txList, args.Error(1)
}
