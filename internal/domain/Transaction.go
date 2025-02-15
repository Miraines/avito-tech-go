package domain

import "time"

// TransactionType represents the type of a transaction.
type TransactionType string

const (
	// Transfer indicates a coin transfer between users.
	Transfer TransactionType = "transfer"
	// Purchase indicates a purchase transaction from the shop.
	Purchase TransactionType = "purchase"
)

// Transaction represents a coin transaction in the system.
// swagger:model Transaction
type Transaction struct {
	ID         uint            `gorm:"primaryKey"`
	FromUserID uint            `gorm:"not null;index"`
	ToUserID   *uint           `gorm:"index"`
	Amount     int             `gorm:"not null"`
	Type       TransactionType `gorm:"size:20;not null"`
	CreatedAt  time.Time
}
