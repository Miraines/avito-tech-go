package server

import (
	"avito-tech-go/internal/config"
	"avito-tech-go/internal/domain"
	"avito-tech-go/pkg/database"
	"fmt"
	"github.com/gin-gonic/gin"
)

func Run(cfg *config.Config) error {
	db, err := database.NewDBConnection(cfg)
	if err != nil {
		return fmt.Errorf("failed to init db: %w", err)
	}

	if err := db.AutoMigrate(
		&domain.User{},
		&domain.MerchItem{},
		&domain.InventoryItem{},
		&domain.Transaction{}); err != nil {
		return fmt.Errorf("failed to migrate db: %w", err)
	}

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// TODO: Добавить ручки

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	return r.Run(addr)
}
