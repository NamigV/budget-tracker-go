package config

import (
	"fmt"
	"os"
)

type Config struct {
	Addr  string
	DB    DatabaseConfig
	Redis RedisConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

func (db DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		db.Host, db.Port, db.User, db.Password, db.Name, db.SSLMode,
	)
}

func Load() Config {
	return Config{
		Addr: getenv("APP_ADDR", ":8080"),
		DB: DatabaseConfig{
			Host:     getenv("DB_HOST", "localhost"),
			Port:     getenv("DB_PORT", "5433"),
			User:     getenv("DB_USER", "budget_tracker"),
			Password: getenv("DB_PASSWORD", "budget_tracker"),
			Name:     getenv("DB_NAME", "budget_tracker"),
			SSLMode:  getenv("DB_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Addr:     getenv("REDIS_ADDR", "localhost:6380"),
			Password: getenv("REDIS_PASSWORD", ""),
			DB:       0,
		},
	}
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
