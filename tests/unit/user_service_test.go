package unit

import (
	"avito-tech-go/tests/unit/mocks"
	"errors"
	"testing"

	"avito-tech-go/internal/domain"
	"avito-tech-go/internal/services"

	"github.com/stretchr/testify/assert"
)

func TestUserService_GetInfo(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockInvRepo := new(mocks.MockInventoryRepository)
	mockTxRepo := new(mocks.MockTransactionRepository)

	userSvc := services.NewUserService(mockUserRepo, mockInvRepo, mockTxRepo)

	t.Run("user not found", func(t *testing.T) {
		mockUserRepo.On("GetUserByID", uint(999)).
			Return((*domain.User)(nil), nil).Once()

		info, err := userSvc.GetInfo(999)
		assert.Error(t, err)
		assert.Nil(t, info)
		assert.Contains(t, err.Error(), "user not found")

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("db error on get user", func(t *testing.T) {
		mockUserRepo.ExpectedCalls = nil
		mockUserRepo.On("GetUserByID", uint(1)).
			Return((*domain.User)(nil), errors.New("db error")).Once()

		info, err := userSvc.GetInfo(1)
		assert.Error(t, err)
		assert.Nil(t, info)
		assert.Contains(t, err.Error(), "db error")

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("success get info", func(t *testing.T) {
		mockUserRepo.ExpectedCalls = nil
		mockInvRepo.ExpectedCalls = nil
		mockTxRepo.ExpectedCalls = nil

		// Ожидается вызов GetUserByID для основного пользователя (ID=1)
		mockUserRepo.On("GetUserByID", uint(1)).
			Return(&domain.User{ID: 1, Username: "alex", Coins: 1000}, nil).Once()

		mockInvRepo.On("GetAllByUser", uint(1)).
			Return([]domain.InventoryItem{
				{ID: 10, UserID: 1, ItemType: "t-shirt", Quantity: 2},
				{ID: 11, UserID: 1, ItemType: "cup", Quantity: 1},
			}, nil).Once()

		mockTxRepo.On("GetUserTransactions", uint(1)).
			Return([]domain.Transaction{
				// Purchase: отправлено в магазин
				{ID: 1, FromUserID: 1, ToUserID: nil, Amount: 80, Type: domain.Purchase},
				// Transfer: от alex (ID=1) к bob (ID=2)
				{ID: 2, FromUserID: 1, ToUserID: func() *uint { v := uint(2); return &v }(), Amount: 100, Type: domain.Transfer},
				// Transfer: от bob (ID=2) к alex (ID=1)
				{ID: 3, FromUserID: 2, ToUserID: func() *uint { v := uint(1); return &v }(), Amount: 50, Type: domain.Transfer},
			}, nil).Once()

		// Ожидаем вызов GetUsernamesByIDs с параметром []uint{2}
		mockUserRepo.On("GetUsernamesByIDs", []uint{2}).
			Return(map[uint]string{2: "bob"}, nil).Once()

		// Удаляем ожидания для GetUserByID(uint(2)), так как они больше не требуются

		info, err := userSvc.GetInfo(1)
		assert.NoError(t, err)
		assert.NotNil(t, info)
		assert.Equal(t, 1000, info.Coins)
		assert.Len(t, info.Inventory, 2) // ожидаем два предмета: t-shirt и cup

		// Проверка истории транзакций:
		// 1) Purchase -> sent: toUser="shop", amount=80
		// 2) Transfer (from=1, to=2) -> sent: toUser="bob", amount=100
		// 3) Transfer (from=2, to=1) -> received: fromUser="bob", amount=50
		assert.Len(t, info.CoinHistory.Sent, 2)
		assert.Len(t, info.CoinHistory.Received, 1)

		mockUserRepo.AssertExpectations(t)
		mockInvRepo.AssertExpectations(t)
		mockTxRepo.AssertExpectations(t)
	})
}
