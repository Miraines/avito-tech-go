package integration

import (
	"avito-tech-go/internal/domain"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

func setupIntegrationDB(t *testing.T) *gorm.DB {
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory sqlite database: %v", err)
	}
	err = db.AutoMigrate(&domain.User{}, &domain.MerchItem{}, &domain.InventoryItem{}, &domain.Transaction{})
	if err != nil {
		t.Fatalf("failed to migrate models: %v", err)
	}
	return db
}
