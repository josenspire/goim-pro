package main

import (
	"bufio"
	"flag"
	"fmt"
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

	for {
		reader := bufio.NewReader(os.Stdin)
		char, _, _ := reader.ReadRune()
		switch char {
		case 'q':
			logger.Infoln("server is starting to disconnect...")
			server.GracefulStopGRPCServer()
			logger.Infoln("server has been gracefully disconnected!")
		default:
			fmt.Println("server continue to listen...")
		}
	}
}
