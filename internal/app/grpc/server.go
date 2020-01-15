/*
	This is the grpc server for `demo app`(module)
	Include some of the server life circle
*/
package grpc

import (
	"context"
	"fmt"
	"github.com/jinzhu/gorm"
	example "goim-pro/api/protos/example"
	protos "goim-pro/api/protos/salty"
	"goim-pro/config"
	"goim-pro/internal/app/repos/user"
	"goim-pro/internal/app/services"
	mysqlsrv "goim-pro/pkg/db/mysql"
	redsrv "goim-pro/pkg/db/redis"
	"goim-pro/pkg/logs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"time"
)

type GRPCServer struct {
	grpcServer *grpc.Server
}

var logger = logs.GetLogger("INFO")

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
			logger.Errorf("mysql tables initialization fail: %s", err)
		}
	}
	redisDB := redsrv.NewRedisConnection()
	if err := redisDB.Connect(); err != nil {
		logger.Errorf("redis connect error: %v", err.Error())
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
	var interceptor grpc.UnaryServerInterceptor
	interceptor = func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		logger.Info(info)
		return handler(ctx, req)
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
	if !db.HasTable(user.User{}) {
		err = db.Set(
			"gorm:table_options",
			"ENGINE=InnoDB DEFAULT CHARSET=utf8",
		).CreateTable(
			user.User{},
		).Error
		if err != nil {
			logger.Errorf("initial mysql tables [users] error: %v", err)
			return
		}
	}
	return
}

func handleServiceRegister(srv *grpc.Server) {
	var gprcService = services.NewService()
	example.RegisterWaiterServer(srv, gprcService.WaiterServer)
	protos.RegisterSMSServiceServer(srv, gprcService.SMSServer)
	protos.RegisterUserServiceServer(srv, gprcService.UserServer)
}
