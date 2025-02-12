package database

import (
	"fmt"
	"gorm.io/driver/postgres"

	"avito-tech-go/internal/config"
	"gorm.io/gorm"
)

func NewDBConnection(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.DbHost,
		cfg.DbUser,
		cfg.DbPass,
		cfg.DbName,
		cfg.DbPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
