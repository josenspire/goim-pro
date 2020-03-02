/*
	This is the grpc server for `demo app`(module)
	Include some of the server life circle
*/
package grpc

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/jinzhu/gorm"
	demo "goim-pro/api/protos/example"
	protos "goim-pro/api/protos/salty"
	"goim-pro/config"
	"goim-pro/internal/app/repos/user"
	"goim-pro/internal/app/services"
	mysqlsrv "goim-pro/pkg/db/mysql"
	redsrv "goim-pro/pkg/db/redis"
	"goim-pro/pkg/logs"
	"goim-pro/pkg/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"time"
)

type GRPCServer struct {
	grpcServer *grpc.Server
}

var (
	logger  = logs.GetLogger("INFO")
	OpenTLS = true

	myRedis *redis.Client
)
// server constructor
func NewServer() *GRPCServer {
	return &GRPCServer{}
}

// initialize server config and db
func (gs *GRPCServer) InitServer() {
	mysqlDB := mysqlsrv.NewMysqlConnection()
	if err := mysqlDB.Connect(); err != nil {
		logger.Errorf("mysql connect error: %v", err)
	} else {
		if err := initialMysqlTables(mysqlDB.GetMysqlInstance()); err != nil {
			logger.Fatalf("mysql tables initialization fail: %s", err)
		}
	}
	rdsClient := redsrv.NewRedisConnection()
	if err := rdsClient.Connect(); err != nil {
		logger.Errorf("redis connect error: %v", err.Error())
	} else {
		myRedis = rdsClient.GetRedisClient()
	}
}

// create gprc server connection
func (gs *GRPCServer) ConnectGRPCServer() {
	tcpAddress := fmt.Sprintf("%s:%s", config.GetAppHost(), config.GetAppPort())
	lis, err := net.Listen("tcp", tcpAddress)
	if err != nil {
		logger.Fatalf("[GRPC server] - [%s] startup failed: %v", tcpAddress, err)
	} else {
		logger.Infof("[GRPC server] - [%s] started successfully!", tcpAddress)
	}

	var opts []grpc.ServerOption

	keepaliveParams := grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle:     time.Second * 60,
		MaxConnectionAge:      time.Hour * 2,
		MaxConnectionAgeGrace: time.Second * 20,
		Time:                  time.Second * 60,
		Timeout:               time.Second * 20,
	})
	opts = append(opts, keepaliveParams)
	// gRPC 拦截器
	//var interceptor grpc.UnaryServerInterceptor
	//interceptor = func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	//	logger.Info(req)
	//	var pb proto.Message
	//	if err := utils.UnmarshalGRPCReq(req.(*protos.GrpcReq), pb); err != nil {
	//		logger.Error(err.Error())
	//	}
	//	return handler(ctx, pb)
	//}
	var interceptor grpc.UnaryServerInterceptor
	interceptor = func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		logger.Info(req, info.FullMethod)
		gRPCReq := req.(*protos.GrpcReq)
		token := gRPCReq.GetToken()

		// handle method on white list
		logger.Infof("[request method]: %s", info.FullMethod)
		if isOnWhiteList(info.FullMethod) {
			return handler(ctx, req)
		}

		// TODO: for local demo
		if token == "1234567890" {
			return handler(ctx, req)
		}
		if utils.IsEmptyStrings(token) {
			resp, _ = utils.NewGRPCResp(http.StatusNonAuthoritativeInfo, nil, "unauthorized access to this resource")
			return resp, nil
		} else {
			// TODO: maybe remove token verify logic and only query from redis
			isValid, payload, err := utils.TokenVerify(token)
			logger.Infof("[userID]: %s", string(payload))
			if err != nil {
				resp, _ = utils.NewGRPCResp(http.StatusInternalServerError, nil, err.Error())
				return resp, nil
			}
			if !isValid {
				resp, _ = utils.NewGRPCResp(http.StatusUnauthorized, nil, "token validation failed")
				return resp, nil
			}

			redisToken, err := myRedis.Get(fmt.Sprintf("TK-%s", string(payload))).Bytes()
			if redisToken == nil {
				resp, _ = utils.NewGRPCResp(http.StatusUnauthorized, nil, "the token has expired")
				return resp, nil
			}
			if err != nil {
				logger.Errorf("interceptor redis read error: %v", err)
				resp, _ = utils.NewGRPCResp(http.StatusInternalServerError, nil, err.Error())
				return resp, nil
			}
			gRPCReq.Token = string(payload)
			return handler(ctx, req)
		}
	}

	opts = append(opts, grpc.UnaryInterceptor(interceptor))

	// 创建 gRPC 服务
	srv := grpc.NewServer(opts...)

	// TODO: should add the service discovery
	// 注册接口服务
	handleServiceRegister(srv)

	// 在 gRPC 服务器上注册反射服务
	reflection.Register(srv)

	// unblock
	go func() {
		// 将监听交给 gRPC 服务处理
		if err = srv.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	gs.grpcServer = srv
}

// stop grpc server by graceful
func (gs *GRPCServer) GracefulStopGRPCServer() {
	gs.grpcServer.GracefulStop()
}

// stop grpc server by force
func (gs *GRPCServer) ForceStopGRPCServer() {
	gs.grpcServer.Stop()
}

func initialMysqlTables(db *gorm.DB) (err error) {
	//if !db.HasTable(&wuid.Wuid{}) {
	//	err = db.Set(
	//		"gorm:table_options",
	//		"ENGINE=InnoDB DEFAULT CHARSET=utf8",
	//	).CreateTable(&wuid.Wuid{}).Error
	//	if err != nil {
	//		logger.Errorf("initial mysql tables [wuids] error: %s", err.Error())
	//		return
	//	}
	//}

	if !db.HasTable(&user.User{}) {
		err = db.Set(
			"gorm:table_options",
			"ENGINE=InnoDB DEFAULT CHARSET=utf8",
		).CreateTable(&user.User{}).Error
		if err != nil {
			logger.Errorf("initial mysql tables [users] error: %s", err.Error())
			return
		}
	}
	return
}

func handleServiceRegister(srv *grpc.Server) {
	var gprcService = services.NewService()
	demo.RegisterWaiterServer(srv, gprcService.WaiterServer)
	protos.RegisterSMSServiceServer(srv, gprcService.SMSServer)
	protos.RegisterUserServiceServer(srv, gprcService.UserServer)
}
