package domain

// MerchItem represents a merchandise item available for purchase.
// swagger:model MerchItem
type MerchItem struct {
	Price    int    `gorm:"not null"`
	ItemType string `gorm:"primaryKey;not null;size:100"`
}
