package services

import (
	"avito-tech-go/internal/domain"
	"avito-tech-go/internal/repositories"
	"fmt"
	"gorm.io/gorm"
)

type MerchService interface {
	BuyItem(userID uint, itemType string) error
}

type merchService struct {
	merchRepo repositories.MerchRepository
	userRepo  repositories.UserRepository
	txRepo    repositories.TransactionRepository
	invRepo   repositories.InventoryRepository
	db        *gorm.DB
}

func NewMerchService(
	merchRepo repositories.MerchRepository,
	userRepo repositories.UserRepository,
	txRepo repositories.TransactionRepository,
	invRepo repositories.InventoryRepository,
	db *gorm.DB,
) MerchService {
	return &merchService{
		merchRepo: merchRepo,
		userRepo:  userRepo,
		txRepo:    txRepo,
		invRepo:   invRepo,
		db:        db}
}

func (m *merchService) BuyItem(userID uint, itemType string) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		merchItem, err := m.merchRepo.GetMerchItemByType(itemType)
		if err != nil {
			return err
		}
		if merchItem == nil {
			return fmt.Errorf("merch item '%s' not found", itemType)
		}

		user, err := m.userRepo.GetUserByID(userID)
		if err != nil {
			return err
		}
		if user == nil {
			return fmt.Errorf("user %d not found", userID)
		}

		if user.Coins < merchItem.Price {
			return fmt.Errorf("user %d does not have enough coins", userID)
		}

		user.Coins -= merchItem.Price
		if err := m.userRepo.UpdateUser(user); err != nil {
			return err
		}

		invItem, err := m.invRepo.GetByUserAndType(userID, itemType)
		if err != nil {
			return err
		}
		if invItem == nil {
			invItem = &domain.InventoryItem{
				ItemType: itemType,
				UserID:   userID,
				Quantity: 1,
			}
			if err := m.invRepo.CreateItem(invItem); err != nil {
				return err
			}
		} else {
			invItem.Quantity++
			if err := m.invRepo.UpdateItem(invItem); err != nil {
				return err
			}
		}

		txItem := &domain.Transaction{
			FromUserID: userID,
			Amount:     merchItem.Price,
			Type:       domain.Purchase,
			ToUserID:   nil,
		}
		if err := m.txRepo.CreateTransaction(txItem); err != nil {
			return err
		}

		return nil
	})
}
