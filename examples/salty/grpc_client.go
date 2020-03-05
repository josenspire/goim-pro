package main

import (
	"fmt"
	protos "goim-pro/api/protos/salty"
	"goim-pro/pkg/logs"
	"google.golang.org/grpc"
	"log"
)

const (
	// address = "111.231.238.209:9090"
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
			case "s1":
				// create Writer service's client
				t := protos.NewSMSServiceClient(conn)
				obtainSMSCode(t, protos.ObtainSMSCodeReq_REGISTER)
			case "s2":
				// create Writer service's client
				t := protos.NewSMSServiceClient(conn)
				obtainSMSCode(t, protos.ObtainSMSCodeReq_LOGIN)
			case "s3":
				t := protos.NewSMSServiceClient(conn)
				obtainSMSCode(t, protos.ObtainSMSCodeReq_RESET_PASSWORD)
			case "rst1":
				t := protos.NewUserServiceClient(conn)
				resetPasswordByTelephone(t, "OLD_PASSWORD")
			case "rst2":
				t := protos.NewUserServiceClient(conn)
				resetPasswordByTelephone(t, "VERIFICATION")
			case "r":
				t := protos.NewUserServiceClient(conn)
				register(t)
			case "lt":
				t := protos.NewUserServiceClient(conn)
				login(t, "TELEPHONE")
			case "le":
				t := protos.NewUserServiceClient(conn)
				login(t, "EMAIL")
			case "lq":
				t := protos.NewUserServiceClient(conn)
				logout(t)
			case "gu":
				t := protos.NewUserServiceClient(conn)
				getUserInfo(t)
			case "qt":
				t := protos.NewUserServiceClient(conn)
				queryUserInfo(t, "TELEPHONE")
			case "qe":
				t := protos.NewUserServiceClient(conn)
				queryUserInfo(t, "EMAIL")
			case "ud":
				t := protos.NewUserServiceClient(conn)
				updateUserInfo(t)

			// contacts
			case "ct-rq":
				t := protos.NewContactServiceClient(conn)
				requestContact(t)
			case "q":
				logger.Infoln("grpc client disconnected!")
				exitChain <- str
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
	logger.Info("** ['s1']: obtainSMSCode - register**")
	logger.Info("** ['s2']: obtainSMSCode - login**")
	logger.Info("** ['s3']: obtainSMSCode - resetPassword**")
	logger.Info("** ['rst1']: resetPassword by telephone with oldPassword**")
	logger.Info("** ['rst2']: resetPassword by telephone with verification**")
	logger.Info("** ['r']: register **")
	logger.Info("** ['lt']: login by telephone **")
	logger.Info("** ['le']: login by email **")
	logger.Info("** ['lq']: user logout **")
	logger.Info("** ['gu']: get user info by userId **")
	logger.Info("** ['qt']: query user info by telephone **")
	logger.Info("** ['qe']: query user info by email **")
	logger.Info("** ['ud']: update user profile **")

	logger.Info("** *********************** **")
	logger.Info("** ['ct-rq']: request add contact **")

	logger.Info("** ['q']: exist [GRPC] client **")
}
