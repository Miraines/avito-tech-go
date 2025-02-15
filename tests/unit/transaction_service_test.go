package unit

import (
	"fmt"
	"testing"

	"avito-tech-go/internal/domain"
	"avito-tech-go/internal/repositories"
	"avito-tech-go/internal/services"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory sqlite database: %v", err)
	}
	err = db.AutoMigrate(&domain.User{}, &domain.Transaction{})
	if err != nil {
		t.Fatalf("failed to migrate models: %v", err)
	}
	return db
}

func TestTransactionService_TransferCoins(t *testing.T) {
	t.Run("same user", func(t *testing.T) {
		db := setupTestDB(t)
		userRepo := repositories.NewUserRepository(db)
		txRepo := repositories.NewTransactionRepository(db)
		txService := services.NewTransactionService(userRepo, txRepo, db)

		err := txService.TransferCoins(1, 1, 100)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot transfer coins to yourself")
	})

	t.Run("amount <= 0", func(t *testing.T) {
		db := setupTestDB(t)
		userRepo := repositories.NewUserRepository(db)
		txRepo := repositories.NewTransactionRepository(db)
		txService := services.NewTransactionService(userRepo, txRepo, db)

		err := txService.TransferCoins(1, 2, 0)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must be greater than 0")
	})

	t.Run("from user not found", func(t *testing.T) {
		db := setupTestDB(t)
		userRepo := repositories.NewUserRepository(db)
		txRepo := repositories.NewTransactionRepository(db)
		txService := services.NewTransactionService(userRepo, txRepo, db)

		err := txService.TransferCoins(1, 2, 10)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user 1 not found")
	})

	t.Run("to user not found", func(t *testing.T) {
		db := setupTestDB(t)
		userRepo := repositories.NewUserRepository(db)
		txRepo := repositories.NewTransactionRepository(db)
		user1 := &domain.User{Coins: 100}
		err := userRepo.CreateUser(user1)
		if err != nil {
			t.Fatalf("failed to create user1: %v", err)
		}
		txService := services.NewTransactionService(userRepo, txRepo, db)

		err = txService.TransferCoins(user1.ID, 9999, 10)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), fmt.Sprintf("user %d not found", 9999))
	})

	t.Run("not enough coins", func(t *testing.T) {
		db := setupTestDB(t)
		userRepo := repositories.NewUserRepository(db)
		txRepo := repositories.NewTransactionRepository(db)
		user1 := &domain.User{Username: "user1", Coins: 5}
		user2 := &domain.User{Username: "user2", Coins: 50}
		if err := userRepo.CreateUser(user1); err != nil {
			t.Fatalf("failed to create user1: %v", err)
		}
		if err := userRepo.CreateUser(user2); err != nil {
			t.Fatalf("failed to create user2: %v", err)
		}
		txService := services.NewTransactionService(userRepo, txRepo, db)

		err := txService.TransferCoins(user1.ID, user2.ID, 10)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), fmt.Sprintf("user %d does not have enough coins", user1.ID))
	})

	t.Run("success", func(t *testing.T) {
		db := setupTestDB(t)
		userRepo := repositories.NewUserRepository(db)
		txRepo := repositories.NewTransactionRepository(db)
		user1 := &domain.User{Username: "user1", Coins: 100}
		user2 := &domain.User{Username: "user2", Coins: 50}
		if err := userRepo.CreateUser(user1); err != nil {
			t.Fatalf("failed to create user1: %v", err)
		}
		if err := userRepo.CreateUser(user2); err != nil {
			t.Fatalf("failed to create user2: %v", err)
		}
		txService := services.NewTransactionService(userRepo, txRepo, db)

		err := txService.TransferCoins(user1.ID, user2.ID, 10)
		assert.NoError(t, err)

		updatedUser1, err := userRepo.GetUserByID(user1.ID)
		if err != nil {
			t.Fatalf("failed to get user1: %v", err)
		}
		updatedUser2, err := userRepo.GetUserByID(user2.ID)
		if err != nil {
			t.Fatalf("failed to get user2: %v", err)
		}
		assert.Equal(t, 90, updatedUser1.Coins)
		assert.Equal(t, 60, updatedUser2.Coins)

		txs, err := txRepo.GetUserTransactions(user1.ID)
		if err != nil {
			t.Fatalf("failed to get transactions: %v", err)
		}
		assert.Len(t, txs, 1)
		txRecord := txs[0]
		assert.Equal(t, user1.ID, txRecord.FromUserID)
		assert.NotNil(t, txRecord.ToUserID)
		assert.Equal(t, user2.ID, *txRecord.ToUserID)
		assert.Equal(t, 10, txRecord.Amount)
		assert.Equal(t, domain.Transfer, txRecord.Type)
	})
}
