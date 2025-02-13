package services

import (
	"avito-tech-go/internal/domain"
	"avito-tech-go/internal/repositories"
	"fmt"
)

type InfoResponse struct {
	Coins       int         `json:"coins"`
	Inventory   []ItemInfo  `json:"inventory"`
	CoinHistory CoinHistory `json:"coinHistory"`
}

type ItemInfo struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

type CoinHistory struct {
	Received []ReceivedTransaction `json:"received"`
	Sent     []SentTransaction     `json:"sent"`
}

type ReceivedTransaction struct {
	FromUser string `json:"fromUser"`
	Amount   int    `json:"amount"`
}

type SentTransaction struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

type UserService interface {
	GetInfo(userID uint) (*InfoResponse, error)
}

type userService struct {
	userRepo        repositories.UserRepository
	invRepo         repositories.InventoryRepository
	transactionRepo repositories.TransactionRepository
}

func NewUserService(
	userRepo repositories.UserRepository,
	invRepo repositories.InventoryRepository,
	txRepo repositories.TransactionRepository,
) UserService {
	return &userService{
		userRepo:        userRepo,
		invRepo:         invRepo,
		transactionRepo: txRepo,
	}
}

func (s *userService) GetInfo(userID uint) (*InfoResponse, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	invItems, err := s.invRepo.GetAllByUser(userID)
	if err != nil {
		return nil, err
	}

	transactions, err := s.transactionRepo.GetUserTransactions(userID)
	if err != nil {
		return nil, err
	}

	info := &InfoResponse{
		Coins:       user.Coins,
		Inventory:   []ItemInfo{},
		CoinHistory: CoinHistory{},
	}

	for _, item := range invItems {
		info.Inventory = append(info.Inventory, ItemInfo{
			Type:     item.ItemType,
			Quantity: item.Quantity,
		})
	}

	for _, tx := range transactions {
		switch tx.Type {
		case domain.Transfer:
			if tx.FromUserID == userID {
				toName := s.getUsernameByID(*tx.ToUserID)
				info.CoinHistory.Sent = append(info.CoinHistory.Sent, SentTransaction{
					ToUser: toName,
					Amount: tx.Amount,
				})
			} else {
				fromName := s.getUsernameByID(tx.FromUserID)
				info.CoinHistory.Received = append(info.CoinHistory.Received, ReceivedTransaction{
					FromUser: fromName,
					Amount:   tx.Amount,
				})
			}

		case domain.Purchase:
			info.CoinHistory.Sent = append(info.CoinHistory.Sent, SentTransaction{
				ToUser: "shop",
				Amount: tx.Amount,
			})

		default:
		}
	}

	return info, nil
}

func (s *userService) getUsernameByID(userID uint) string {
	u, err := s.userRepo.GetUserByID(userID)
	if err != nil || u == nil {
		return fmt.Sprintf("user_%d", userID)
	}
	return u.Username
}
