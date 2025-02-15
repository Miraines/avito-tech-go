package database

import (
	"avito-tech-go/internal/domain"
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"time"

	"avito-tech-go/internal/config"
	"gorm.io/gorm"
)

func NewDBConnection(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBName,
		cfg.DBPort)

	var db *gorm.DB
	var err error
	maxAttempts := 10
	for i := 1; i <= maxAttempts; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			sqlDB, err := db.DB()
			if err == nil && sqlDB.Ping() == nil {
				// Настраиваем пул соединений
				sqlDB.SetMaxOpenConns(95)                  // Максимальное число открытых соединений
				sqlDB.SetMaxIdleConns(50)                  // Максимальное число простаивающих соединений
				sqlDB.SetConnMaxLifetime(15 * time.Minute) // Время жизни соединения

				return db, nil
			}
		}
		fmt.Printf("Попытка подключения к базе (%d/%d) не удалась: %v\n", i, maxAttempts, err)
		time.Sleep(5 * time.Second)
	}
	return nil, errors.New("database not ready after multiple attempts")
}

func SeedMerch(db *gorm.DB) error {
	items := []domain.MerchItem{
		{ItemType: "t-shirt", Price: 80},
		{ItemType: "cup", Price: 20},
		{ItemType: "book", Price: 50},
		{ItemType: "pen", Price: 10},
		{ItemType: "powerbank", Price: 200},
		{ItemType: "hoody", Price: 300},
		{ItemType: "umbrella", Price: 200},
		{ItemType: "socks", Price: 10},
		{ItemType: "wallet", Price: 50},
		{ItemType: "pink-hoody", Price: 500},
	}

	for _, item := range items {
		if err := db.Where("item_type = ?", item.ItemType).
			FirstOrCreate(&item).Error; err != nil {
			return err
		}
	}
	return nil
}
