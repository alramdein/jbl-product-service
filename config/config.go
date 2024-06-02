package config

import (
	"errors"
	"os"
	"referral-system/util"
	"time"

	"github.com/labstack/gommon/log"
)

type Config struct {
	Db              DbConfig
	JwtSecret       string
	ReferralLinkExp time.Time
}

type DbConfig struct {
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassword string
	DbName     string
	DbSSLMode  string
}

func GetConfig() (*Config, error) {
	db := DbConfig{
		DbHost:     os.Getenv("DB_HOST"),
		DbPort:     os.Getenv("DB_PORT"),
		DbUser:     os.Getenv("DB_USER"),
		DbPassword: os.Getenv("DB_PASSWORD"),
		DbName:     os.Getenv("DB_NAME"),
		DbSSLMode:  os.Getenv("DB_SSLMODE"),
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, errors.New("jwt secret required")
	}

	now := time.Now().UTC()
	str := os.Getenv("REFERRAL_LINK_EXP")
	refferalLinkExp, err := util.ParseDurationString(str)
	if err != nil {
		log.Error(err)
		return nil, errors.New("failed to get config jwt referral link expiration")
	}
	if refferalLinkExp.IsZero() {
		refferalLinkExp = now.Add(DefaultReferralLinkExp)
	}

	return &Config{
		Db:              db,
		JwtSecret:       jwtSecret,
		ReferralLinkExp: refferalLinkExp,
	}, nil
}
