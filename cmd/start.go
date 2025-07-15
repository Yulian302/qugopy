package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Yulian302/qugopy/config"
	"github.com/Yulian302/qugopy/grpc"
	"github.com/Yulian302/qugopy/logging"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/spf13/cobra"

	"github.com/Yulian302/qugopy/shell"
)

var startCmd *cobra.Command

func RunApp(isProduction bool) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	cfg, err := config.LoadConfig()

	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// process debug mode params
	if !isProduction {
		envMode := os.Getenv("MODE")
		if envMode == "" {
			envMode = "local"
		}
		config.AppConfig.MODE = envMode

		envWorkers := os.Getenv("WORKERS")
		if envWorkers != "" {
			if parsed, err := strconv.Atoi(envWorkers); err == nil {
				config.AppConfig.WORKERS = parsed
			} else {
				fmt.Println("Invalid WORKERS value:", err)
			}
		} else {
			config.AppConfig.WORKERS = 2
		}
	} else {
		gin.SetMode(gin.ReleaseMode)
		f, _ := os.Create(fmt.Sprintf(filepath.Join(config.ProjectRootPath, "gin.log")))
		gin.DefaultWriter = io.MultiWriter(f)
	}

	// process cli params
	if modeFlag := startCmd.Flag("mode").Value.String(); modeFlag != "" && startCmd.Flag("mode").Changed {
		cfg.MODE = modeFlag
	}
	if workersFlag := startCmd.Flag("workers").Value.String(); workersFlag != "" && startCmd.Flag("workers").Changed {
		if workers, err := strconv.Atoi(workersFlag); err == nil {
			cfg.WORKERS = workers
		}
	}

	// set up redis
	if config.AppConfig.MODE == "redis" {
		rdb = redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%s", config.AppConfig.REDIS.HOST, config.AppConfig.REDIS.PORT),
		})
		logging.DebugLog(fmt.Sprintf("Successfully connected to Redis (host: %s, port: %s)", config.AppConfig.REDIS.HOST, config.AppConfig.REDIS.PORT))
	}

	errCh := make(chan error, 2)

	if cfg.MODE == "local" {
		go func() { errCh <- grpc.Start() }()
		time.Sleep(100 * time.Millisecond)
	}

	var cancel context.CancelFunc
	cancel, err = StartApp(cfg.MODE, cfg.WORKERS, isProduction)
	if err != nil {
		log.Fatalf("App failed to start: %v", err)
	}

	// start shell only in prod
	if isProduction {
		shell.StartInteractiveShell(rdb)
	}

	select {
	case err := <-errCh:
		fmt.Println("Service exited with error:", err)
	case <-ctx.Done():
		fmt.Println("Shutting down gracefully...")
		cancel()
	}

	stop()
}

func init() {
	startCmd = &cobra.Command{
		Use:   "start",
		Short: "Start the application",
		Run: func(cmd *cobra.Command, args []string) {
			RunApp(true)
		},
	}

	startCmd.Flags().StringP("mode", "m", "local", "mode for queuing tasks: redis | local")
	startCmd.Flags().IntP("workers", "w", 2, "number of concurrent workers")
	rootCmd.AddCommand(startCmd)
}
