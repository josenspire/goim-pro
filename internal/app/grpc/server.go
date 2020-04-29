/*
	This is the grpc server for `demo app`(module)
	Include some of the server life circle
*/
package grpc

import (
	"context"
	"fmt"
	"github.com/jinzhu/gorm"
	demo "goim-pro/api/protos/example"
	protos "goim-pro/api/protos/salty"
	"goim-pro/config"
	"goim-pro/internal/app/models"
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

	myRedis *redsrv.BaseClient
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
func (gs *GRPCServer) StartGRPCServer() {
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
		if isOnWhiteList(info.FullMethod) {
			return handler(ctx, req)
		}

		// TODO: for local demo
		if token == "1234567890" {
			gRPCReq.Token = "01E07SG858N3CGV5M1APVQKZYR"
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

			redisToken := myRedis.Get(fmt.Sprintf("TK-%s", string(payload)))
			if redisToken == "" {
				resp, _ = utils.NewGRPCResp(http.StatusUnauthorized, nil, "the token has expired")
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
	logger.Info("graceful shutdown, waiting for all processes done...")
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

	if !db.HasTable(&models.User{}) {
		err = db.Set(
			"gorm:table_options",
			"ENGINE=InnoDB DEFAULT CHARSET=utf8",
		).CreateTable(&models.User{}).Error
		if err != nil {
			logger.Errorf("initial mysql tables [users] error: %s", err.Error())
			return
		}
	}
	if !db.HasTable(&models.Contact{}) {
		err = db.Set(
			"gorm:table_options",
			"ENGINE=InnoDB DEFAULT CHARSET=utf8",
		).CreateTable(&models.Contact{}).Error
		if err != nil {
			logger.Errorf("initial mysql tables [contacts] error: %s", err.Error())
			return
		}

		err = db.Model(&models.Contact{}).AddForeignKey("userId", "users(userId)", "CASCADE", "CASCADE").Error
		if err != nil {
			logger.Errorf("init table constraint relation error: %s", err.Error())
		}
	}
	return
}

func handleServiceRegister(srv *grpc.Server) {
	var grpcService = services.NewService()
	demo.RegisterWaiterServer(srv, grpcService.WaiterServer)
	protos.RegisterSMSServiceServer(srv, grpcService.SMSServer)
	protos.RegisterUserServiceServer(srv, grpcService.UserServer)
	protos.RegisterContactServiceServer(srv, grpcService.ContactServer)
	protos.RegisterGroupServiceServer(srv, grpcService.GroupServer)
}
