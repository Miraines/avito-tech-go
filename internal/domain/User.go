package domain

import "time"

// User represents an employee in the system.
// swagger:model User
type User struct {
	ID           uint   `gorm:"primaryKey"`
	Username     string `gorm:"uniqueIndex;not null;size:255"`
	PasswordHash string `gorm:"not null; size:255"`
	Coins        int    `gorm:"default:1000"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
