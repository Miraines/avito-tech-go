package config

import (
	"os"
	"strconv"
)

type Config struct {
	AppPort   string
	DBHost    string
	DBPort    int
	DBUser    string
	DBPass    string
	DBName    string
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
		DBHost:    getEnv("DB_HOST", "localhost"),
		DBPort:    dbPort,
		DBUser:    getEnv("DB_USER", "postgres"),
		DBPass:    getEnv("DB_PASSWORD", "postgres"),
		DBName:    getEnv("DB_NAME", "avito_shop"),
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
