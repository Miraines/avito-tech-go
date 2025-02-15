package integration

import (
	"testing"

	"avito-tech-go/internal/domain"
	"avito-tech-go/internal/repositories"
	"avito-tech-go/internal/services"
	"github.com/stretchr/testify/assert"
)

func TestIntegration_MerchPurchase_NewItem(t *testing.T) {
	db := setupIntegrationDB(t)

	merchRepo := repositories.NewMerchRepository(db)
	userRepo := repositories.NewUserRepository(db)
	invRepo := repositories.NewInventoryRepository(db)
	txRepo := repositories.NewTransactionRepository(db)

	merchItem := &domain.MerchItem{
		ItemType: "t-shirt",
		Price:    80,
	}
	err := merchRepo.CreateMerchItem(merchItem)
	assert.NoError(t, err)

	user := &domain.User{
		Username:     "testuser",
		PasswordHash: "irrelevant",
		Coins:        100,
	}
	err = userRepo.CreateUser(user)
	assert.NoError(t, err)

	merchService := services.NewMerchService(merchRepo, userRepo, txRepo, invRepo, db)

	err = merchService.BuyItem(user.ID, "t-shirt")
	assert.NoError(t, err)

	updatedUser, err := userRepo.GetUserByID(user.ID)
	assert.NoError(t, err)
	expectedCoins := user.Coins - merchItem.Price
	assert.Equal(t, expectedCoins, updatedUser.Coins)

	invItem, err := invRepo.GetByUserAndType(user.ID, "t-shirt")
	assert.NoError(t, err)
	assert.NotNil(t, invItem)
	assert.Equal(t, 1, invItem.Quantity)

	txs, err := txRepo.GetUserTransactions(user.ID)
	assert.NoError(t, err)
	assert.Len(t, txs, 1)
	txRecord := txs[0]
	assert.Equal(t, user.ID, txRecord.FromUserID)
	assert.Nil(t, txRecord.ToUserID)
	assert.Equal(t, merchItem.Price, txRecord.Amount)
	assert.Equal(t, domain.Purchase, txRecord.Type)
}

func TestIntegration_MerchPurchase_ExistingItem(t *testing.T) {
	db := setupIntegrationDB(t)

	merchRepo := repositories.NewMerchRepository(db)
	userRepo := repositories.NewUserRepository(db)
	invRepo := repositories.NewInventoryRepository(db)
	txRepo := repositories.NewTransactionRepository(db)

	merchItem := &domain.MerchItem{
		ItemType: "t-shirt",
		Price:    80,
	}
	err := merchRepo.CreateMerchItem(merchItem)
	assert.NoError(t, err)

	user := &domain.User{
		Username:     "testuser2",
		PasswordHash: "irrelevant",
		Coins:        300,
	}
	err = userRepo.CreateUser(user)
	assert.NoError(t, err)

	invItem := &domain.InventoryItem{
		UserID:   user.ID,
		ItemType: "t-shirt",
		Quantity: 1,
	}
	err = invRepo.CreateItem(invItem)
	assert.NoError(t, err)

	merchService := services.NewMerchService(merchRepo, userRepo, txRepo, invRepo, db)

	err = merchService.BuyItem(user.ID, "t-shirt")
	assert.NoError(t, err)

	updatedUser, err := userRepo.GetUserByID(user.ID)
	assert.NoError(t, err)
	expectedCoins := user.Coins - merchItem.Price
	assert.Equal(t, expectedCoins, updatedUser.Coins)

	updatedInvItem, err := invRepo.GetByUserAndType(user.ID, "t-shirt")
	assert.NoError(t, err)
	assert.NotNil(t, updatedInvItem)
	assert.Equal(t, 2, updatedInvItem.Quantity)

	txs, err := txRepo.GetUserTransactions(user.ID)
	assert.NoError(t, err)
	assert.Len(t, txs, 1)
	txRecord := txs[0]
	assert.Equal(t, user.ID, txRecord.FromUserID)
	assert.Nil(t, txRecord.ToUserID)
	assert.Equal(t, merchItem.Price, txRecord.Amount)
	assert.Equal(t, domain.Purchase, txRecord.Type)
}

func TestIntegration_MerchPurchase_InsufficientFunds(t *testing.T) {
	db := setupIntegrationDB(t)

	merchRepo := repositories.NewMerchRepository(db)
	userRepo := repositories.NewUserRepository(db)
	invRepo := repositories.NewInventoryRepository(db)
	txRepo := repositories.NewTransactionRepository(db)

	merchItem := &domain.MerchItem{
		ItemType: "t-shirt",
		Price:    1500,
	}
	err := merchRepo.CreateMerchItem(merchItem)
	assert.NoError(t, err)

	user := &domain.User{
		Username:     "insufficient",
		PasswordHash: "irrelevant",
	}
	err = userRepo.CreateUser(user)
	assert.NoError(t, err)

	merchService := services.NewMerchService(merchRepo, userRepo, txRepo, invRepo, db)
	err = merchService.BuyItem(user.ID, "t-shirt")
	assert.Error(t, err)

	updatedUser, err := userRepo.GetUserByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, 1000, updatedUser.Coins)

	invItem, err := invRepo.GetByUserAndType(user.ID, "t-shirt")
	assert.NoError(t, err)
	assert.Nil(t, invItem)

	txs, err := txRepo.GetUserTransactions(user.ID)
	assert.NoError(t, err)
	assert.Len(t, txs, 0)
}

func TestIntegration_MerchPurchase_InvalidItem(t *testing.T) {
	db := setupIntegrationDB(t)

	merchRepo := repositories.NewMerchRepository(db)
	userRepo := repositories.NewUserRepository(db)
	invRepo := repositories.NewInventoryRepository(db)
	txRepo := repositories.NewTransactionRepository(db)

	user := &domain.User{
		Username:     "invaliditem",
		PasswordHash: "irrelevant",
		Coins:        200,
	}
	err := userRepo.CreateUser(user)
	assert.NoError(t, err)

	merchService := services.NewMerchService(merchRepo, userRepo, txRepo, invRepo, db)
	err = merchService.BuyItem(user.ID, "non-existent-item")
	assert.Error(t, err)

	updatedUser, err := userRepo.GetUserByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.Coins, updatedUser.Coins)

	invItem, err := invRepo.GetByUserAndType(user.ID, "non-existent-item")
	assert.NoError(t, err)
	assert.Nil(t, invItem)

	txs, err := txRepo.GetUserTransactions(user.ID)
	assert.NoError(t, err)
	assert.Len(t, txs, 0)
}

func TestIntegration_MerchPurchase_ExactFunds(t *testing.T) {
	db := setupIntegrationDB(t)

	merchRepo := repositories.NewMerchRepository(db)
	userRepo := repositories.NewUserRepository(db)
	invRepo := repositories.NewInventoryRepository(db)
	txRepo := repositories.NewTransactionRepository(db)

	merchItem := &domain.MerchItem{
		ItemType: "t-shirt",
		Price:    80,
	}
	err := merchRepo.CreateMerchItem(merchItem)
	assert.NoError(t, err)

	user := &domain.User{
		Username:     "exactfunds",
		PasswordHash: "irrelevant",
		Coins:        80,
	}
	err = userRepo.CreateUser(user)
	assert.NoError(t, err)

	merchService := services.NewMerchService(merchRepo, userRepo, txRepo, invRepo, db)
	err = merchService.BuyItem(user.ID, "t-shirt")
	assert.NoError(t, err)

	updatedUser, err := userRepo.GetUserByID(user.ID)
	assert.NoError(t, err)
	expectedCoins := user.Coins - merchItem.Price
	assert.Equal(t, expectedCoins, updatedUser.Coins)

	invItem, err := invRepo.GetByUserAndType(user.ID, "t-shirt")
	assert.NoError(t, err)
	assert.NotNil(t, invItem)
	assert.Equal(t, 1, invItem.Quantity)

	txs, err := txRepo.GetUserTransactions(user.ID)
	assert.NoError(t, err)
	assert.Len(t, txs, 1)
	txRecord := txs[0]
	assert.Equal(t, user.ID, txRecord.FromUserID)
	assert.Nil(t, txRecord.ToUserID)
	assert.Equal(t, merchItem.Price, txRecord.Amount)
	assert.Equal(t, domain.Purchase, txRecord.Type)
}
