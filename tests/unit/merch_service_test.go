package unit

import (
	"testing"

	"avito-tech-go/internal/domain"
	"avito-tech-go/internal/services"
	"avito-tech-go/tests/unit/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestMerchService_BuyItem(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory sqlite database: %v", err)
	}

	err = db.AutoMigrate(&domain.User{}, &domain.MerchItem{}, &domain.InventoryItem{}, &domain.Transaction{})
	if err != nil {
		t.Fatalf("failed to migrate models: %v", err)
	}

	mockMerchRepo := new(mocks.MockMerchRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockTxRepo := new(mocks.MockTransactionRepository)
	mockInvRepo := new(mocks.MockInventoryRepository)

	merchSvc := services.NewMerchService(mockMerchRepo, mockUserRepo, mockTxRepo, mockInvRepo, db)

	t.Run("merch item not found", func(t *testing.T) {
		mockMerchRepo.ExpectedCalls = nil

		mockMerchRepo.On("GetMerchItemByType", "unknown").
			Return((*domain.MerchItem)(nil), nil).Once()

		err := merchSvc.BuyItem(1, "unknown")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "merch item 'unknown' not found")

		mockMerchRepo.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		mockMerchRepo.ExpectedCalls = nil
		mockUserRepo.ExpectedCalls = nil

		mockMerchRepo.On("GetMerchItemByType", "t-shirt").
			Return(&domain.MerchItem{ItemType: "t-shirt", Price: 80}, nil).Once()

		mockUserRepo.On("GetUserByID", uint(10)).
			Return((*domain.User)(nil), nil).Once()

		err := merchSvc.BuyItem(10, "t-shirt")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user 10 not found")

		mockMerchRepo.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("not enough coins", func(t *testing.T) {
		mockMerchRepo.ExpectedCalls = nil
		mockUserRepo.ExpectedCalls = nil

		mockMerchRepo.On("GetMerchItemByType", "t-shirt").
			Return(&domain.MerchItem{ItemType: "t-shirt", Price: 80}, nil).Once()

		mockUserRepo.On("GetUserByID", uint(1)).
			Return(&domain.User{ID: 1, Coins: 50}, nil).Once()

		err := merchSvc.BuyItem(1, "t-shirt")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "does not have enough coins")

		mockMerchRepo.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("success buy (new item)", func(t *testing.T) {
		mockMerchRepo.ExpectedCalls = nil
		mockUserRepo.ExpectedCalls = nil
		mockInvRepo.ExpectedCalls = nil
		mockTxRepo.ExpectedCalls = nil

		mockMerchRepo.On("GetMerchItemByType", "t-shirt").
			Return(&domain.MerchItem{ItemType: "t-shirt", Price: 80}, nil).Once()

		mockUserRepo.On("GetUserByID", uint(1)).
			Return(&domain.User{ID: 1, Coins: 100}, nil).Once()

		mockUserRepo.On("UpdateUser", mock.MatchedBy(func(user *domain.User) bool {
			return user.ID == 1 && user.Coins == 20
		})).Return(nil).Once()

		mockInvRepo.On("GetByUserAndType", uint(1), "t-shirt").
			Return((*domain.InventoryItem)(nil), nil).Once()

		mockInvRepo.On("CreateItem", mock.MatchedBy(func(item *domain.InventoryItem) bool {
			return item.UserID == 1 && item.ItemType == "t-shirt" && item.Quantity == 1
		})).Return(nil).Once()

		mockTxRepo.On("CreateTransaction", mock.MatchedBy(func(tx *domain.Transaction) bool {
			return tx.FromUserID == 1 && tx.Amount == 80 && tx.Type == domain.Purchase
		})).Return(nil).Once()

		err := merchSvc.BuyItem(1, "t-shirt")
		assert.NoError(t, err)

		mockMerchRepo.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
		mockInvRepo.AssertExpectations(t)
		mockTxRepo.AssertExpectations(t)
	})

	t.Run("success buy (existing item)", func(t *testing.T) {
		mockMerchRepo.ExpectedCalls = nil
		mockUserRepo.ExpectedCalls = nil
		mockInvRepo.ExpectedCalls = nil
		mockTxRepo.ExpectedCalls = nil

		mockMerchRepo.On("GetMerchItemByType", "t-shirt").
			Return(&domain.MerchItem{ItemType: "t-shirt", Price: 80}, nil).Once()

		mockUserRepo.On("GetUserByID", uint(2)).
			Return(&domain.User{ID: 2, Coins: 300}, nil).Once()

		mockUserRepo.On("UpdateUser", mock.MatchedBy(func(user *domain.User) bool {
			return user.ID == 2 && user.Coins == 220
		})).Return(nil).Once()

		mockInvRepo.On("GetByUserAndType", uint(2), "t-shirt").
			Return(&domain.InventoryItem{ID: 100, Quantity: 1, UserID: 2, ItemType: "t-shirt"}, nil).Once()

		mockInvRepo.On("UpdateItem", mock.MatchedBy(func(item *domain.InventoryItem) bool {
			return item.ID == 100 && item.Quantity == 2
		})).Return(nil).Once()

		mockTxRepo.On("CreateTransaction", mock.MatchedBy(func(tx *domain.Transaction) bool {
			return tx.FromUserID == 2 && tx.Amount == 80 && tx.Type == domain.Purchase
		})).Return(nil).Once()

		err := merchSvc.BuyItem(2, "t-shirt")
		assert.NoError(t, err)

		mockMerchRepo.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
		mockInvRepo.AssertExpectations(t)
		mockTxRepo.AssertExpectations(t)
	})
}
