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
	userRepo        repositories.UserRepository
	transactionRepo repositories.TransactionRepository
	db              *gorm.DB
}

func NewTransactionService(
	userRepo repositories.UserRepository,
	txRepo repositories.TransactionRepository,
	db *gorm.DB,
) TransactionService {
	return &transactionService{userRepo: userRepo, transactionRepo: txRepo, db: db}
}

func (t *transactionService) TransferCoins(fromUserID, toUserID uint, amount int) error {
	if fromUserID == toUserID {
		return fmt.Errorf("cannot transfer coins to yourself")
	}

	if amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}

	err := t.db.Transaction(func(tx *gorm.DB) error {
		userRepoTx := repositories.NewUserRepository(tx)
		txRepoTx := repositories.NewTransactionRepository(tx)

		fromUser, err := userRepoTx.GetUserByID(fromUserID)
		if err != nil {
			return err
		}
		if fromUser == nil {
			return fmt.Errorf("user %d not found", fromUserID)
		}

		toUser, err := userRepoTx.GetUserByID(toUserID)
		if err != nil {
			return err
		}
		if toUser == nil {
			return fmt.Errorf("user %d not found", toUserID)
		}

		if fromUser.Coins < amount {
			return fmt.Errorf("user %d does not have enough coins", fromUserID)
		}

		fromUser.Coins -= amount
		toUser.Coins += amount

		if err := userRepoTx.UpdateUser(fromUser); err != nil {
			return err
		}
		if err := userRepoTx.UpdateUser(toUser); err != nil {
			return err
		}

		transaction := &domain.Transaction{
			FromUserID: fromUserID,
			ToUserID:   &toUserID,
			Amount:     amount,
			Type:       domain.Transfer,
		}

		if err := txRepoTx.CreateTransaction(transaction); err != nil {
			return err
		}

		return nil
	})

	return err
}
