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
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

var logger = logs.GetLogger("INFO")

func New() *grpc.Server {
	// TODO: should add the service discovery
	tcpAddress := fmt.Sprintf("%s:%s", configs.GetAppHost(), configs.GetAppPort())
	lis, err := net.Listen("tcp", tcpAddress)
	if err != nil {
		logger.Fatalf("server - [%s] startup failed: %v\n", err)
	} else {
		logger.Infof("server - [%s] started successfully...", tcpAddress)
	}

	// 创建 gRPC 服务
	srv := grpc.NewServer()
	// 注册接口服务
	var userService = services.NewService()
	protos.RegisterWaiterServer(srv, userService.WaiterServer)
	protos.RegisterUserServiceServer(srv, userService.UserServer)
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
