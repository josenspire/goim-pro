package main

import (
	"fmt"
	protos "goim-pro/api/protos/salty"
	"goim-pro/pkg/logs"
	"google.golang.org/grpc"
	"log"
)

const (
	//address = "111.231.238.209:9090"
	address = "127.0.0.1:9090"
)

var logger = logs.GetLogger("INFO")

func main() {
	//interceptor := grpc.WithUnaryInterceptor(func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	//	logger.Info(req)
	//	pb := req.(proto.Message)
	//	logger.Info(pb)
	//	//gprcReq := protos.GrpcReq{
	//	//	DeviceID: "",
	//	//	Version:  "",
	//	//	Language: 0,
	//	//	Os:       0,
	//	//	Token:    "",
	//	//	Data:     *any.Any{},
	//	//}
	//	return nil
	//})

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("grpc connect fail: %v", err)
	}
	defer conn.Close()

	exitChain := make(chan string)
	var str string
	go func() {
		for {
			_, _ = fmt.Scanln(&str)
			switch str {
			case "s":
				// create Writer service's client
				t := protos.NewSMSServiceClient(conn)
				obtainSMSCode(t)
				break
			case "r":
				t := protos.NewUserServiceClient(conn)
				register(t)
				break
			case "t":
				t := protos.NewUserServiceClient(conn)
				login(t, "TELEPHONE")
				break
			case "e":
				t := protos.NewUserServiceClient(conn)
				login(t, "EMAIL")
				break
			case "q":
				logger.Infoln("grpc client disconnected!")
				exitChain <- str
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
	logger.Info("**** Welcome to [GRPC] client tools ****")
	logger.Info("**** Can input the commons to test ****")
	logger.Info("** ['s']: obtainSMSCode **")
	logger.Info("** ['r']: register **")
	logger.Info("** ['l1']: login by telephone **")
	logger.Info("** ['l2']: login by email **")

	logger.Info("** ['q']: exist [GRPC] client **")
}
