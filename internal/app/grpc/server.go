/*
	This is the grpc server for `demo app`(module)
	Include some of the server life circle
*/
package grpc

import (
	"fmt"
	"goim-pro/api/protos"
	"goim-pro/config"
	"goim-pro/internal/app/services"
	"goim-pro/pkg/db/mysql"
	"goim-pro/pkg/db/redis"
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

func init() {
	mysqlDB := mysqlsrv.NewMysqlConnection()
	if err := mysqlDB.Connect(); err != nil {
		logger.Panicf("mysql connect error: %v", err)
	}
	if err := redsrv.NewRedisService().Connect(); err != nil {
		logger.Panicf("redis connect error: %v", err)
	}
}

func NewServer() *GRPCServer {
	return &GRPCServer{}
}
func (gs *GRPCServer) ConnectGRPCServer() {
	tcpAddress := fmt.Sprintf("%s:%s", config.GetAppHost(), config.GetAppPort())
	lis, err := net.Listen("tcp", tcpAddress)
	if err != nil {
		logger.Fatalf("[GRPC server] - [%s] startup failed: %v\n", err)
	} else {
		logger.Infof("[GRPC server] - [%s] started successfully!\n", tcpAddress)
	}

	keepaliveParams := grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle:     time.Second * 60,
		MaxConnectionAge:      time.Hour * 2,
		MaxConnectionAgeGrace: time.Second * 20,
		Time:                  time.Second * 60,
		Timeout:               time.Second * 20,
	})
	// 创建 gRPC 服务
	srv := grpc.NewServer(keepaliveParams)

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

func (gs *GRPCServer) GracefulStopGRPCServer() {
	gs.grpcServer.GracefulStop()
}

func (gs *GRPCServer) ForceStopGRPCServer() {
	gs.grpcServer.Stop()
}

func handleServiceRegister(srv *grpc.Server) {
	var userService = services.NewService()
	protos.RegisterWaiterServer(srv, userService.WaiterServer)
	protos.RegisterUserServiceServer(srv, userService.UserServer)
}
