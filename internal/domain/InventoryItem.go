package domain

import "time"

// InventoryItem represents an item in a user's inventory.
// swagger:model InventoryItem
type InventoryItem struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"not null;index;foreignKey:UserID;references:ID"`
	Quantity  int    `gorm:"default:0"`
	ItemType  string `gorm:"not null;size:100"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
