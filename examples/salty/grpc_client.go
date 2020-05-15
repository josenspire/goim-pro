package main

import (
	"fmt"
	protos "goim-pro/api/protos/salty"
	"goim-pro/examples/salty/contact"
	"goim-pro/examples/salty/group"
	"goim-pro/examples/salty/user"
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
				for i := 1; i <= 3; i++ {
					user.ObtainSMSCode(t, protos.ObtainSMSCodeReq_REGISTER, i)
				}
			case "s2":
				// create Writer service's client
				t := protos.NewSMSServiceClient(conn)
				user.ObtainSMSCode(t, protos.ObtainSMSCodeReq_LOGIN, 3)
			case "s3":
				t := protos.NewSMSServiceClient(conn)
				user.ObtainSMSCode(t, protos.ObtainSMSCodeReq_RESET_PASSWORD, 3)
			case "rst1":
				t := protos.NewUserServiceClient(conn)
				user.ResetPasswordByTelephone(t, "OLD_PASSWORD")
			case "rst2":
				t := protos.NewUserServiceClient(conn)
				user.ResetPasswordByTelephone(t, "VERIFICATION")
			case "r":
				t := protos.NewUserServiceClient(conn)
				user.Register(t)
			case "lt":
				t := protos.NewUserServiceClient(conn)
				for i := 1; i <= 3; i++ {
					user.Login(t, "TELEPHONE", i)
				}
			case "lt2":
				t := protos.NewUserServiceClient(conn)
				user.Login(t, "TELEPHONE", 3)
			case "lt3":
				t := protos.NewUserServiceClient(conn)
				user.LoginWithCode(t, "TELEPHONE")
			case "le":
				t := protos.NewUserServiceClient(conn)
				for i := 1; i <= 3; i++ {
					user.Login(t, "EMAIL", i)
				}
			case "lq":
				t := protos.NewUserServiceClient(conn)
				user.Logout(t)
			case "gu":
				t := protos.NewUserServiceClient(conn)
				user.GetUserInfo(t)
			case "qt":
				t := protos.NewUserServiceClient(conn)
				user.QueryUserInfo(t, "TELEPHONE")
			case "qe":
				t := protos.NewUserServiceClient(conn)
				user.QueryUserInfo(t, "EMAIL")
			case "ud":
				t := protos.NewUserServiceClient(conn)
				user.UpdateUserInfo(t)

			// contacts
			case "ct-rq":
				t := protos.NewContactServiceClient(conn)
				contact.RequestContact(t)
			case "ct-acp":
				t := protos.NewContactServiceClient(conn)
				contact.AcceptContact(t)
			case "ct-rf":
				t := protos.NewContactServiceClient(conn)
				contact.RefusedContact(t)
			case "ct-udt":
				t := protos.NewContactServiceClient(conn)
				contact.UpdateRemarkInfo(t)
			case "ct-fa":
				t := protos.NewContactServiceClient(conn)
				contact.GetContacts(t)

			// groups
			case "gp-ct":
				t := protos.NewGroupServiceClient(conn)
				group.CreateGroup(t)
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
	logger.Info("** ['lt2']: login by telephone with diff deviceId **")
	logger.Info("** [lt3]: login by telephone with sms code **")
	logger.Info("** ['le']: login by email **")
	logger.Info("** ['lq']: user logout **")
	logger.Info("** ['gu']: get user info by userId **")
	logger.Info("** ['qt']: query user info by telephone **")
	logger.Info("** ['qe']: query user info by email **")
	logger.Info("** ['ud']: update user profile **")

	logger.Info("** *********************** **")
	logger.Info("** ['ct-rq']: request add contact **")
	logger.Info("** ['ct-acp']: accept add contact **")
	logger.Info("** ['ct-rf']: refused contact request **")
	logger.Info("** ['ct-udt']: update contact remark profile **")
	logger.Info("** ['ct-fa']: get all contacts **")

	logger.Info("** *********************** **")
	logger.Info("** ['gp-ct']: create new group **")

	logger.Info("** ['q']: exist [GRPC] client **")
}
