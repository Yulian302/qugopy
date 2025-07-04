package cmd

import (
	"github.com/Yulian302/qugopy/grpc"
	"github.com/spf13/cobra"
)

var (
	mode    string
	workers uint8
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the application",
	Run: func(cmd *cobra.Command, args []string) {
		go grpc.Start()
		StartApp(mode)
	},
}

func init() {
	startCmd.PersistentFlags().StringVarP(&mode, "mode", "m", "local", "mode for queuing tasks: redis | local")
	startCmd.PersistentFlags().Uint8VarP(&workers, "workers", "w", 2, "number of concurrent workers")
	rootCmd.AddCommand(startCmd)
}
