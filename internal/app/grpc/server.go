/*
	This is the grpc server for `demo app`(module)
	Include some of the server life circle
*/
package grpc

import (
	"fmt"
	protos "goim-pro/api/protos"
	"goim-pro/configs"
	"goim-pro/internal/app/services"
	"goim-pro/pkg/logs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"time"
)

var logger = logs.GetLogger("INFO")

func New() *grpc.Server {
	tcpAddress := fmt.Sprintf("%s:%s", configs.GetAppHost(), configs.GetAppPort())
	lis, err := net.Listen("tcp", tcpAddress)
	if err != nil {
		logger.Fatalf("server - [%s] startup failed: %v\n", err)
	} else {
		logger.Infof("server - [%s] started successfully...", tcpAddress)
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

	return srv
}

func handleServiceRegister(srv *grpc.Server) {
	var userService = services.NewService()
	protos.RegisterWaiterServer(srv, userService.WaiterServer)
	protos.RegisterUserServiceServer(srv, userService.UserServer)
}
