package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/Yulian302/qugopy/config"
	"github.com/Yulian302/qugopy/internal/api"
	w "github.com/Yulian302/qugopy/workers"
	"github.com/go-redis/redis"
)

var rdb *redis.Client

func StartApp(mode string, workers int, isProduction bool) (context.CancelFunc, error) {
	wd := w.NewWorkerDistributor(rdb)
	cancel, err := wd.DistributeWorkers(workers, mode, isProduction, rdb)
	if err != nil {
		return nil, fmt.Errorf("failed to distribute workers: %w", err)
	}

	router := api.NewRouter(rdb)

	go func() {
		if err := router.Run(fmt.Sprintf("%s:%s", config.AppConfig.HOST, config.AppConfig.PORT)); err != nil {
			log.Printf("failed to start server: %v", err)
		}
	}()

	return cancel, nil
}
