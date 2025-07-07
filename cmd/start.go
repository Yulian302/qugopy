package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/Yulian302/qugopy/grpc"
	"github.com/spf13/cobra"

	"github.com/Yulian302/qugopy/shell"
)

var (
	mode    string
	workers uint8
)

func RunApp(isProduction bool) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	errCh := make(chan error, 2)

	go func() { errCh <- StartApp(mode, isProduction) }()
	go func() { errCh <- grpc.Start(isProduction) }()

	if isProduction {
		shell.StartInteractiveShell()
	}

	select {
	case err := <-errCh:
		fmt.Println("Service exited with error:", err)
	case <-ctx.Done():
		fmt.Println("Shutting down gracefully...")
	}

	stop()
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the application",
	Run: func(cmd *cobra.Command, args []string) {
		RunApp(true)
	},
}

func init() {
	startCmd.PersistentFlags().StringVarP(&mode, "mode", "m", "local", "mode for queuing tasks: redis | local")
	startCmd.PersistentFlags().Uint8VarP(&workers, "workers", "w", 2, "number of concurrent workers")
	rootCmd.AddCommand(startCmd)
}
