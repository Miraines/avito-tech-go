package config

import (
	"os"
	"strconv"
)

type Config struct {
	AppPort   string
	DbHost    string
	DbPort    int
	DbUser    string
	DbPass    string
	DbName    string
	JWTSecret string
}

func LoadConfig() (*Config, error) {
	dbPortStr := os.Getenv("DB_PORT")
	if dbPortStr == "" {
		dbPortStr = "5432"
	}
	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		AppPort:   getEnv("APP_PORT", "8080"),
		DbHost:    getEnv("DB_HOST", "localhost"),
		DbPort:    dbPort,
		DbUser:    getEnv("DB_USER", "postgres"),
		DbPass:    getEnv("DB_PASSWORD", "postgres"),
		DbName:    getEnv("DB_NAME", "avito_shop"),
		JWTSecret: getEnv("JWT_SECRET", "avitomiraines"),
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	return val
}
