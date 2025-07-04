package main

import (
	"os"

	"github.com/Yulian302/qugopy/cmd"
)

func main() {
	runMode := os.Getenv("RUN_MODE")
	if runMode == "air" {
		cmd.RunDev()
	} else {
		cmd.Execute()
	}
}
