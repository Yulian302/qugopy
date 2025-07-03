package config

import (
	"errors"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type RedisConfig struct {
	HOST string
	PORT string
}

type AppConfig struct {
	HOST  string
	PORT  string
	REDIS RedisConfig
}

func LoadConfig() (*AppConfig, error) {

	cfg := &AppConfig{
		HOST: os.Getenv("HOST"),
		PORT: os.Getenv("PORT"),
		REDIS: RedisConfig{
			HOST: os.Getenv("REDIS_HOST"),
			PORT: os.Getenv("REDIS_PORT"),
		},
	}

	if cfg.HOST == "" || cfg.PORT == "" {
		return nil, errors.New("configuration error")
	}

	return cfg, nil
}
