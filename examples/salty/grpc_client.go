package main

import (
	"bufio"
	protos "goim-pro/api/protos/salty"
	"goim-pro/pkg/logs"
	"google.golang.org/grpc"
	"log"
	"os"
)

const (
	//address = "111.231.238.209:9090"
	address = "127.0.0.1:9090"
)

var logger = logs.GetLogger("INFO")

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("grpc connect fail: %v", err)
	}
	defer conn.Close()

	exitChain := make(chan string)
	go func() {
		for {
			reader := bufio.NewReader(os.Stdin)
			char, _, _ := reader.ReadRune()
			switch char {
			case 's':
				// create Writer service's client
				t := protos.NewSMSServiceClient(conn)
				obtainSMSCode(t)
				break
			case 'r':
				t := protos.NewUserServiceClient(conn)
				register(t)
				break
			case 'q':
				logger.Infoln("grpc client disconnected!")
				exitChain <- string(char)
				break
			default:
				logger.Info("server continue to listen...")
			}
			logger.Info("********************************************")
		}
	}()

	toolsIntroduce()

	_ = <-exitChain
	logger.Info("grpc server exit!")
}

func toolsIntroduce() {
	logger.Info("********************************************")
	logger.Info("**** welcome to grpc client tools ****")
	logger.Info("**** can input the commons to test ****")
	logger.Info("** 's': obtainSMSCode **")
	logger.Info("** 'r': register **")

	logger.Info("** 'q': exist grpc client **")
}
