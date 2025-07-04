package cmd

import "github.com/Yulian302/qugopy/grpc"

func RunDev() {
	go grpc.Start()
	StartApp(mode)
}
