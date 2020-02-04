package main

import (
	"bufio"
	"errors"
	"fmt"
	"goim-pro/internal/app/grpc"
	"goim-pro/pkg/logs"
	"os"
)

var logger = logs.GetLogger("INFO")
var server *grpc.GRPCServer

func main() {
	server = grpc.NewServer()
	server.InitServer()
	server.ConnectGRPCServer()

	defer func() {
		if r := recover(); r != nil {
			err := errors.New(fmt.Sprint(r))
			logger.Panicf("uncheck exception: %v", err)
		}
	}()

	exitChain := make(chan string)
	go func() {
		for {
			reader := bufio.NewReader(os.Stdin)
			char, _, _ := reader.ReadRune()
			str := string(char)
			switch str {
			case "q":
				logger.Infoln("server is starting to disconnect...")
				server.GracefulStopGRPCServer()
				logger.Infoln("server has been gracefully disconnected!")
				exitChain <- str
			default:
				logger.Info("server continue to listen...")
			}
		}
	}()
	str := <-exitChain
	logger.Info("exit!", str)
}
