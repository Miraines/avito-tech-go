package main

import (
	"log"

	"avito-tech-go/internal/config"
	"avito-tech-go/internal/server"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	log.Printf("Конфигурация успешно загружена: %+v", cfg)

	if err := server.Run(cfg); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
