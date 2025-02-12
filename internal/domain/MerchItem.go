package domain

type MerchItem struct {
	Price    int    `gorm:"not null"`
	ItemType string `gorm:"primaryKey;not null;size:100"`
}
