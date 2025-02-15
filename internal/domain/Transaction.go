package domain

import "time"

type TransactionType string

const (
	Transfer TransactionType = "transfer"
	Purchase TransactionType = "purchase"
)

type Transaction struct {
	ID         uint            `gorm:"primaryKey"`
	FromUserID uint            `gorm:"not null;index"` // При покупке: from = пользователь
	ToUserID   *uint           `gorm:"index"`          // при переводе: ToUser = другой user; при покупке: NULL
	Amount     int             `gorm:"not null"`
	Type       TransactionType `gorm:"size:20;not null"`
	CreatedAt  time.Time
}
