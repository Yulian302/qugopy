package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/Yulian302/qugopy/config"
	"github.com/Yulian302/qugopy/internal/api"
	"github.com/go-redis/redis"
)

var rdb *redis.Client
var ctx = context.Background()

func StartApp(mode string) {
	cfg, err := config.LoadConfig()

	if err != nil {
		log.Fatal(err)
	}
	if config.Mode(mode).IsValid() {
		cfg.MODE = config.Mode(mode)
	} else {
		cfg.MODE = "local"
	}

	rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.REDIS.HOST, cfg.REDIS.PORT),
	})

	router := api.NewRouter(rdb)

	if err := router.Run(fmt.Sprintf("%s:%s", cfg.HOST, cfg.PORT)); err != nil {
		log.Fatal("Failed to start server: %w", err)
	}
}
