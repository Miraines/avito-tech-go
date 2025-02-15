package integration

import (
	"testing"

	"avito-tech-go/internal/domain"
	"avito-tech-go/internal/repositories"
	"avito-tech-go/internal/services"
	"github.com/stretchr/testify/assert"
)

func TestIntegration_UserInfo(t *testing.T) {
	t.Run("User not found", func(t *testing.T) {
		db := setupIntegrationDB(t)
		userRepo := repositories.NewUserRepository(db)
		invRepo := repositories.NewInventoryRepository(db)
		txRepo := repositories.NewTransactionRepository(db)
		userService := services.NewUserService(userRepo, invRepo, txRepo)

		info, err := userService.GetInfo(9999) // несуществующий ID
		assert.Error(t, err)
		assert.Nil(t, info)
	})

	t.Run("User with no inventory and no transactions", func(t *testing.T) {
		db := setupIntegrationDB(t)
		userRepo := repositories.NewUserRepository(db)
		invRepo := repositories.NewInventoryRepository(db)
		txRepo := repositories.NewTransactionRepository(db)
		userService := services.NewUserService(userRepo, invRepo, txRepo)

		user := &domain.User{
			Username: "emptyuser",
		}
		err := userRepo.CreateUser(user)
		assert.NoError(t, err)

		info, err := userService.GetInfo(user.ID)
		assert.NoError(t, err)
		assert.Equal(t, 1000, info.Coins)
		assert.Empty(t, info.Inventory)
		assert.Empty(t, info.CoinHistory.Received)
		assert.Empty(t, info.CoinHistory.Sent)
	})

	t.Run("User with inventory", func(t *testing.T) {
		db := setupIntegrationDB(t)
		userRepo := repositories.NewUserRepository(db)
		invRepo := repositories.NewInventoryRepository(db)
		txRepo := repositories.NewTransactionRepository(db)
		userService := services.NewUserService(userRepo, invRepo, txRepo)

		user := &domain.User{Username: "invuser"}
		err := userRepo.CreateUser(user)
		assert.NoError(t, err)

		item1 := &domain.InventoryItem{
			UserID:   user.ID,
			ItemType: "t-shirt",
			Quantity: 2,
		}
		assert.NoError(t, invRepo.CreateItem(item1))

		item2 := &domain.InventoryItem{
			UserID:   user.ID,
			ItemType: "cup",
			Quantity: 3,
		}
		assert.NoError(t, invRepo.CreateItem(item2))

		info, err := userService.GetInfo(user.ID)
		assert.NoError(t, err)
		assert.Len(t, info.Inventory, 2)

		invMap := make(map[string]int)
		for _, inv := range info.Inventory {
			invMap[inv.Type] = inv.Quantity
		}
		assert.Equal(t, 2, invMap["t-shirt"])
		assert.Equal(t, 3, invMap["cup"])
	})

	t.Run("User with purchase transaction", func(t *testing.T) {
		db := setupIntegrationDB(t)
		userRepo := repositories.NewUserRepository(db)
		invRepo := repositories.NewInventoryRepository(db)
		txRepo := repositories.NewTransactionRepository(db)
		userService := services.NewUserService(userRepo, invRepo, txRepo)

		user := &domain.User{Username: "purchaser"}
		assert.NoError(t, userRepo.CreateUser(user))

		purchaseTx := &domain.Transaction{
			FromUserID: user.ID,
			ToUserID:   nil,
			Amount:     80,
			Type:       domain.Purchase,
		}
		assert.NoError(t, txRepo.CreateTransaction(purchaseTx))

		info, err := userService.GetInfo(user.ID)
		assert.NoError(t, err)
		assert.Len(t, info.CoinHistory.Sent, 1)
		sentTx := info.CoinHistory.Sent[0]
		assert.Equal(t, "shop", sentTx.ToUser)
		assert.Equal(t, 80, sentTx.Amount)
	})

	t.Run("User with transfer transaction as sender", func(t *testing.T) {
		db := setupIntegrationDB(t)
		userRepo := repositories.NewUserRepository(db)
		invRepo := repositories.NewInventoryRepository(db)
		txRepo := repositories.NewTransactionRepository(db)
		userService := services.NewUserService(userRepo, invRepo, txRepo)

		sender := &domain.User{Username: "sender"}
		receiver := &domain.User{Username: "receiver"}
		assert.NoError(t, userRepo.CreateUser(sender))
		assert.NoError(t, userRepo.CreateUser(receiver))

		transferTx := &domain.Transaction{
			FromUserID: sender.ID,
			ToUserID:   func(u uint) *uint { return &u }(receiver.ID),
			Amount:     150,
			Type:       domain.Transfer,
		}
		assert.NoError(t, txRepo.CreateTransaction(transferTx))

		info, err := userService.GetInfo(sender.ID)
		assert.NoError(t, err)
		assert.Len(t, info.CoinHistory.Sent, 1)
		sentTx := info.CoinHistory.Sent[0]
		assert.Equal(t, receiver.Username, sentTx.ToUser)
		assert.Equal(t, 150, sentTx.Amount)
	})

	t.Run("User with transfer transaction as receiver", func(t *testing.T) {
		db := setupIntegrationDB(t)
		userRepo := repositories.NewUserRepository(db)
		invRepo := repositories.NewInventoryRepository(db)
		txRepo := repositories.NewTransactionRepository(db)
		userService := services.NewUserService(userRepo, invRepo, txRepo)

		sender := &domain.User{Username: "sender2"}
		receiver := &domain.User{Username: "receiver2"}
		assert.NoError(t, userRepo.CreateUser(sender))
		assert.NoError(t, userRepo.CreateUser(receiver))

		transferTx := &domain.Transaction{
			FromUserID: sender.ID,
			ToUserID:   func(u uint) *uint { return &u }(receiver.ID),
			Amount:     200,
			Type:       domain.Transfer,
		}
		assert.NoError(t, txRepo.CreateTransaction(transferTx))

		info, err := userService.GetInfo(receiver.ID)
		assert.NoError(t, err)
		assert.Len(t, info.CoinHistory.Received, 1)
		receivedTx := info.CoinHistory.Received[0]
		assert.Equal(t, sender.Username, receivedTx.FromUser)
		assert.Equal(t, 200, receivedTx.Amount)
	})

	t.Run("User with mixed inventory and transactions", func(t *testing.T) {
		db := setupIntegrationDB(t)
		userRepo := repositories.NewUserRepository(db)
		invRepo := repositories.NewInventoryRepository(db)
		txRepo := repositories.NewTransactionRepository(db)
		userService := services.NewUserService(userRepo, invRepo, txRepo)

		user := &domain.User{Username: "mixeduser"}
		assert.NoError(t, userRepo.CreateUser(user))

		item := &domain.InventoryItem{
			UserID:   user.ID,
			ItemType: "hat",
			Quantity: 1,
		}
		assert.NoError(t, invRepo.CreateItem(item))

		purchaseTx := &domain.Transaction{
			FromUserID: user.ID,
			ToUserID:   nil,
			Amount:     50,
			Type:       domain.Purchase,
		}
		assert.NoError(t, txRepo.CreateTransaction(purchaseTx))

		receiver := &domain.User{Username: "mixedReceiver"}
		assert.NoError(t, userRepo.CreateUser(receiver))
		transferTx := &domain.Transaction{
			FromUserID: user.ID,
			ToUserID:   func(u uint) *uint { return &u }(receiver.ID),
			Amount:     100,
			Type:       domain.Transfer,
		}
		assert.NoError(t, txRepo.CreateTransaction(transferTx))

		info, err := userService.GetInfo(user.ID)
		assert.NoError(t, err)
		assert.Equal(t, 1000, info.Coins)
		assert.Len(t, info.Inventory, 1)
		assert.Equal(t, "hat", info.Inventory[0].Type)
		assert.Equal(t, 1, info.Inventory[0].Quantity)

		assert.Len(t, info.CoinHistory.Sent, 2)
		sentMap := make(map[string]int)
		for _, st := range info.CoinHistory.Sent {
			sentMap[st.ToUser] = st.Amount
		}
		assert.Equal(t, 50, sentMap["shop"])
		assert.Equal(t, 100, sentMap[receiver.Username])
	})
}
