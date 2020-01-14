package main

import (
	"bufio"
	"flag"
	"goim-pro/internal/app/grpc"
	"goim-pro/pkg/logs"
	"os"
)

var logger = logs.GetLogger("INFO")
var server *grpc.GRPCServer

func main() {
	flag.Parse()

	server = grpc.NewServer()
	server.InitServer()
	server.ConnectGRPCServer()

	reader := bufio.NewReader(os.Stdin)
	for {
		char, _, err := reader.ReadRune()
		if err != nil {
			logger.Errorf("input error: %s", err.Error())
			goto BLOCK
		}
		switch char {
		case 'q':
			logger.Infoln("server is starting to disconnect...")
			server.GracefulStopGRPCServer()
			logger.Infoln("server has been gracefully disconnected!")
			goto BLOCK
		default:
			logger.Info("server continue to listen...")
		}
	}
BLOCK:
	logger.Info("exit!")
}
