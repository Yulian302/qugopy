package cmd

import (
	"fmt"

	"github.com/Yulian302/qugopy/config"
	"github.com/Yulian302/qugopy/internal/api"
	w "github.com/Yulian302/qugopy/workers"
	"github.com/go-redis/redis"
)

var rdb *redis.Client

func StartApp(mode string, workers int, isProduction bool) error {

	wd := w.NewWorkerDistributor(rdb)
	cancel, err := wd.DistributeWorkers(workers, mode, isProduction, rdb)
	if err != nil {
		return fmt.Errorf("failed to distribute workers: %w", err)
	}
	defer cancel()

	router := api.NewRouter(rdb)

	if err := router.Run(fmt.Sprintf("%s:%s", config.AppConfig.HOST, config.AppConfig.PORT)); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}
