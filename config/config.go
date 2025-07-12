package config

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

var (
	// Get current file full path from runtime
	_, b, _, _ = runtime.Caller(0)

	// Root folder of this project
	ProjectRootPath = filepath.Join(filepath.Dir(b), "../")
)

// Mode specified by user. Can be either `redis` or `local`. If `redis` mode is specified, all tasks are pushed to the Redis data store. On the other hand, if `local` is specified, in-memory priority queue is used.
type Mode string

// Checks whether the mode is of valid type. Can be either `redis` or `local`
func (m Mode) IsValid() bool {
	return m == "redis" || m == "local"
}

type RedisConfig struct {
	HOST string
	PORT string
}

type BrevoConfig struct {
	URL     string
	API_KEY string
	EMAIL   string
}

type RootConfig struct {
	HOST    string
	PORT    string
	REDIS   RedisConfig
	BREVO   BrevoConfig
	MODE    string
	WORKERS int
}

func LoadConfig() (*RootConfig, error) {

	if err := godotenv.Load(ProjectRootPath + "/.env"); err != nil {
		log.Fatal("Could not find project .env file")
	}

	cfg := &RootConfig{
		HOST: os.Getenv("HOST"),
		PORT: os.Getenv("PORT"),
		REDIS: RedisConfig{
			HOST: os.Getenv("REDIS_HOST"),
			PORT: os.Getenv("REDIS_PORT"),
		},
		BREVO: BrevoConfig{
			URL:     os.Getenv("BREVO_URL"),
			API_KEY: os.Getenv("BREVO_API_KEY"),
			EMAIL:   os.Getenv("BREVO_EMAIL"),
		},
		MODE:    "local",
		WORKERS: 2,
	}

	if cfg.HOST == "" || cfg.PORT == "" {
		return nil, errors.New("configuration error")
	}
	AppConfig = cfg
	return cfg, nil
}

var AppConfig *RootConfig = &RootConfig{}
