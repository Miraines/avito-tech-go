package services

import (
	"avito-tech-go/internal/domain"
	"avito-tech-go/internal/repositories"
	"fmt"
	"gorm.io/gorm"
)

type TransactionService interface {
	TransferCoins(fromUserID, toUserID uint, amount int) error
}

type transactionService struct {
	transactionRepo repositories.TransactionRepository
	userRepo        repositories.UserRepository
	db              *gorm.DB
}

func NewTransactionServie(
	userRepo repositories.UserRepository,
	tx repositories.TransactionRepository,
	db *gorm.DB,
) TransactionService {
	return &transactionService{userRepo: userRepo, transactionRepo: tx, db: db}
}

func (t *transactionService) TransferCoins(fromUserID, toUserID uint, amount int) error {
	if fromUserID == toUserID {
		return fmt.Errorf("Cannot transfer coins to yourself")
	}

	if amount <= 0 {
		return fmt.Errorf("Amount must be greater than 0")
	}

	err := t.db.Transaction(func(tx *gorm.DB) error {
		fromUser, err := t.userRepo.GetUserByID(fromUserID)
		if err != nil {
			return err
		}
		if fromUser == nil {
			return fmt.Errorf("User %d not found", fromUserID)
		}

		toUser, err := t.userRepo.GetUserByID(toUserID)
		if err != nil {
			return err
		}
		if toUser == nil {
			return fmt.Errorf("User %d not found", toUserID)
		}

		if fromUser.Coins < amount {
			return fmt.Errorf("User %d does not have enough coins", fromUserID)
		}

		fromUser.Coins -= amount
		toUser.Coins += amount

		if err := tx.Save(fromUser).Error; err != nil {
			return err
		}

		if err := tx.Save(toUser).Error; err != nil {
			return err
		}

		transaction := &domain.Transaction{
			FromUserID: fromUserID,
			ToUserID:   &toUserID,
			Amount:     amount,
			Type:       domain.Transfer,
		}

		if err := tx.Create(transaction).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}
