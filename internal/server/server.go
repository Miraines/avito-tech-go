package server

import (
	_ "avito-tech-go/docs"
	"avito-tech-go/internal/config"
	"avito-tech-go/internal/domain"
	"avito-tech-go/internal/handlers"
	"avito-tech-go/internal/middleware"
	"avito-tech-go/internal/repositories"
	"avito-tech-go/internal/services"
	"avito-tech-go/pkg/database"
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           API Avito shop
// @version         1.0.0
// @description     API for the Avito internal merch shop. Employees can purchase merch with coins and transfer coins among themselves.
// @contact.name    API Support
// @contact.email   support@avito.tech
// @host            localhost:8080
// @BasePath        /
// @schemes         http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
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

	if err := database.SeedMerch(db); err != nil {
		return fmt.Errorf("failed to seed merch: %w", err)
	}

	userRepo := repositories.NewUserRepository(db)
	merchRepo := repositories.NewMerchRepository(db)
	invRepo := repositories.NewInventoryRepository(db)
	txRepo := repositories.NewTransactionRepository(db)

	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	userService := services.NewUserService(userRepo, invRepo, txRepo)
	transactionService := services.NewTransactionService(userRepo, txRepo, db)
	merchService := services.NewMerchService(merchRepo, userRepo, txRepo, invRepo, db)

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/api/auth", handlers.AuthHandler(authService))

	authMw := middleware.JWTAuthMiddleware(cfg.JWTSecret)

	r.GET("/api/info", authMw, handlers.InfoHandler(userService))
	r.POST("/api/sendCoin", authMw, handlers.SendCoinHandler(transactionService, userRepo))
	r.GET("/api/buy/:item", authMw, handlers.BuyMerchHandler(merchService))

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	return r.Run(addr)
}
