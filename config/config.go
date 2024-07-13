package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Db        DbConfig
	JwtSecret string
}

type DbConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	db := DbConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	// if jwtSecret == "" {
	// 	return nil, errors.New("jwt secret required")
	// }

	return &Config{
		Db:        db,
		JwtSecret: jwtSecret,
	}, nil
}
