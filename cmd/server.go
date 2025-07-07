package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Yulian302/qugopy/config"
	"github.com/Yulian302/qugopy/internal/api"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

var rdb *redis.Client
var ctx = context.Background()

func StartApp(mode string, isProduction bool) error {

	cfg, err := config.LoadConfig()

	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	if config.Mode(mode).IsValid() {
		cfg.MODE = config.Mode(mode)
	} else {
		cfg.MODE = "local"
	}

	if isProduction {
		gin.SetMode(gin.ReleaseMode)
		f, _ := os.Create(fmt.Sprintf(filepath.Join(config.ProjectRootPath, "gin.log")))
		gin.DefaultWriter = io.MultiWriter(f)
	}

	rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.REDIS.HOST, cfg.REDIS.PORT),
	})

	router := api.NewRouter(rdb)

	if err := router.Run(fmt.Sprintf("%s:%s", cfg.HOST, cfg.PORT)); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}
