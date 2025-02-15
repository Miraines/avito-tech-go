package integration

import (
	"testing"

	"avito-tech-go/internal/domain"
	"avito-tech-go/internal/repositories"
	"avito-tech-go/internal/services"
	"github.com/stretchr/testify/assert"
)

func TestIntegration_CoinTransfer(t *testing.T) {
	db := setupIntegrationDB(t)

	userRepo := repositories.NewUserRepository(db)
	txRepo := repositories.NewTransactionRepository(db)
	transferService := services.NewTransactionService(userRepo, txRepo, db)

	alice := &domain.User{
		Username:     "alice",
		PasswordHash: "irrelevant",
		Coins:        1000,
	}
	err := userRepo.CreateUser(alice)
	assert.NoError(t, err)

	bob := &domain.User{
		Username:     "bob",
		PasswordHash: "irrelevant",
		Coins:        1000,
	}
	err = userRepo.CreateUser(bob)
	assert.NoError(t, err)

	t.Run("Successful transfer", func(t *testing.T) {
		err := transferService.TransferCoins(alice.ID, bob.ID, 200)
		assert.NoError(t, err)

		updatedAlice, err := userRepo.GetUserByID(alice.ID)
		assert.NoError(t, err)
		updatedBob, err := userRepo.GetUserByID(bob.ID)
		assert.NoError(t, err)

		assert.Equal(t, 800, updatedAlice.Coins)
		assert.Equal(t, 1200, updatedBob.Coins)

		txs, err := txRepo.GetUserTransactions(alice.ID)
		assert.NoError(t, err)
		assert.Len(t, txs, 1)
		txRecord := txs[0]
		assert.Equal(t, alice.ID, txRecord.FromUserID)
		assert.NotNil(t, txRecord.ToUserID)
		assert.Equal(t, bob.ID, *txRecord.ToUserID)
		assert.Equal(t, 200, txRecord.Amount)
		assert.Equal(t, domain.Transfer, txRecord.Type)
	})

	t.Run("Insufficient funds", func(t *testing.T) {
		err := transferService.TransferCoins(alice.ID, bob.ID, 900)
		assert.Error(t, err)

		updatedAlice, err := userRepo.GetUserByID(alice.ID)
		assert.NoError(t, err)
		updatedBob, err := userRepo.GetUserByID(bob.ID)
		assert.NoError(t, err)
		assert.Equal(t, 800, updatedAlice.Coins)
		assert.Equal(t, 1200, updatedBob.Coins)
	})

	t.Run("Transfer zero coins", func(t *testing.T) {
		err := transferService.TransferCoins(alice.ID, bob.ID, 0)
		assert.Error(t, err)

		updatedAlice, err := userRepo.GetUserByID(alice.ID)
		assert.NoError(t, err)
		updatedBob, err := userRepo.GetUserByID(bob.ID)
		assert.NoError(t, err)
		assert.Equal(t, 800, updatedAlice.Coins)
		assert.Equal(t, 1200, updatedBob.Coins)
	})

	t.Run("Transfer negative amount", func(t *testing.T) {
		err := transferService.TransferCoins(alice.ID, bob.ID, -50)
		assert.Error(t, err)

		updatedAlice, err := userRepo.GetUserByID(alice.ID)
		assert.NoError(t, err)
		updatedBob, err := userRepo.GetUserByID(bob.ID)
		assert.NoError(t, err)
		assert.Equal(t, 800, updatedAlice.Coins)
		assert.Equal(t, 1200, updatedBob.Coins)
	})

	t.Run("Transfer to non-existing receiver", func(t *testing.T) {
		nonExistingReceiverID := uint(9999)
		err := transferService.TransferCoins(alice.ID, nonExistingReceiverID, 100)
		assert.Error(t, err)

		updatedAlice, err := userRepo.GetUserByID(alice.ID)
		assert.NoError(t, err)
		assert.Equal(t, 800, updatedAlice.Coins)
	})

	t.Run("Transfer from non-existing sender", func(t *testing.T) {
		nonExistingSenderID := uint(9999)
		err := transferService.TransferCoins(nonExistingSenderID, bob.ID, 100)
		assert.Error(t, err)

		updatedBob, err := userRepo.GetUserByID(bob.ID)
		assert.NoError(t, err)
		assert.Equal(t, 1200, updatedBob.Coins)
	})

	t.Run("Transfer to self", func(t *testing.T) {
		err := transferService.TransferCoins(alice.ID, alice.ID, 100)
		assert.Error(t, err)

		updatedAlice, err := userRepo.GetUserByID(alice.ID)
		assert.NoError(t, err)
		assert.Equal(t, 800, updatedAlice.Coins)
	})
}
