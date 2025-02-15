package services

import (
	"avito-tech-go/internal/domain"
	"avito-tech-go/internal/repositories"
	"fmt"
	"golang.org/x/sync/errgroup"
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
	// Сначала получаем данные пользователя
	user, err := s.getUser(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	var (
		inventory    []ItemInfo
		transactions []domain.Transaction
	)

	// Создаем группу для параллельного выполнения
	var g errgroup.Group

	// Запускаем запрос на получение инвентаря
	g.Go(func() error {
		inv, err := s.getInventoryInfo(userID)
		if err != nil {
			return err
		}
		inventory = inv
		return nil
	})

	// Запускаем запрос на получение транзакций
	g.Go(func() error {
		txs, err := s.getUserTransactions(userID)
		if err != nil {
			return err
		}
		transactions = txs
		return nil
	})

	// Ждем завершения обеих горутин
	if err := g.Wait(); err != nil {
		return nil, err
	}

	// Получаем имена пользователей, задействованных в транзакциях
	usernamesMap, err := s.getUsernamesForTransactions(userID, transactions)
	if err != nil {
		return nil, err
	}

	// Формируем историю транзакций
	coinHistory := s.buildCoinHistory(userID, transactions, usernamesMap)

	return &InfoResponse{
		Coins:       user.Coins,
		Inventory:   inventory,
		CoinHistory: coinHistory,
	}, nil
}

// getUser получает пользователя по ID
func (s *userService) getUser(userID uint) (*domain.User, error) {
	return s.userRepo.GetUserByID(userID)
}

// getInventoryInfo получает информацию об инвентаре пользователя
func (s *userService) getInventoryInfo(userID uint) ([]ItemInfo, error) {
	invItems, err := s.invRepo.GetAllByUser(userID)
	if err != nil {
		return nil, err
	}
	var items []ItemInfo
	for _, item := range invItems {
		items = append(items, ItemInfo{
			Type:     item.ItemType,
			Quantity: item.Quantity,
		})
	}
	return items, nil
}

// getUserTransactions получает транзакции пользователя
func (s *userService) getUserTransactions(userID uint) ([]domain.Transaction, error) {
	return s.transactionRepo.GetUserTransactions(userID)
}

// getUsernamesForTransactions собирает уникальные ID и пакетно получает имена пользователей
func (s *userService) getUsernamesForTransactions(userID uint, transactions []domain.Transaction) (map[uint]string, error) {
	uniqueIDs := make(map[uint]struct{})
	for _, tx := range transactions {
		if tx.Type == domain.Transfer {
			if tx.FromUserID != userID {
				uniqueIDs[tx.FromUserID] = struct{}{}
			}
			if tx.ToUserID != nil && *tx.ToUserID != userID {
				uniqueIDs[*tx.ToUserID] = struct{}{}
			}
		}
	}
	var ids []uint
	for id := range uniqueIDs {
		ids = append(ids, id)
	}
	return s.userRepo.GetUsernamesByIDs(ids)
}

// buildCoinHistory формирует историю транзакций для ответа
func (s *userService) buildCoinHistory(userID uint, transactions []domain.Transaction, usernamesMap map[uint]string) CoinHistory {
	var history CoinHistory
	for _, tx := range transactions {
		switch tx.Type {
		case domain.Transfer:
			if tx.FromUserID == userID {
				toName := "unknown"
				if tx.ToUserID != nil {
					if name, ok := usernamesMap[*tx.ToUserID]; ok {
						toName = name
					}
				}
				history.Sent = append(history.Sent, SentTransaction{
					ToUser: toName,
					Amount: tx.Amount,
				})
			} else {
				fromName := "unknown"
				if name, ok := usernamesMap[tx.FromUserID]; ok {
					fromName = name
				}
				history.Received = append(history.Received, ReceivedTransaction{
					FromUser: fromName,
					Amount:   tx.Amount,
				})
			}
		case domain.Purchase:
			history.Sent = append(history.Sent, SentTransaction{
				ToUser: "shop",
				Amount: tx.Amount,
			})
		}
	}
	return history
}

func (s *userService) getUsernameByID(userID uint) string {
	u, err := s.userRepo.GetUserByID(userID)
	if err != nil || u == nil {
		return fmt.Sprintf("user_%d", userID)
	}
	return u.Username
}
