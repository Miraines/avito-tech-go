package repositories

import (
	"avito-tech-go/internal/domain"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreateTransaction(tx *domain.Transaction) error
	GetUserTransactions(userID uint) ([]domain.Transaction, error)
	GetTransactionsByType(userID uint, txType domain.TransactionType) ([]domain.Transaction, error)
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) CreateTransaction(tx *domain.Transaction) error {
	return r.db.Create(tx).Error
}

func (r *transactionRepository) GetUserTransactions(userID uint) ([]domain.Transaction, error) {
	var transactions []domain.Transaction
	err := r.db.
		Where("from_user_id = ? OR to_user_id = ?", userID, userID).
		Find(&transactions).Error
	return transactions, err
}

func (r *transactionRepository) GetTransactionsByType(userID uint, txType domain.TransactionType) ([]domain.Transaction, error) {
	var transactions []domain.Transaction
	err := r.db.Where("(from_user_id = ? OR to_user_id = ?) AND type = ?", userID, userID, txType).Find(&transactions).Error
	return transactions, err
}
