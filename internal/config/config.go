package config

import (
	"fmt"
	"os"
)

type Config struct {
	Addr string
	DB   DatabaseConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
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
	}
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
