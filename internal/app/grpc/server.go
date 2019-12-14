/*
	This is the grpc server for `salty app`(module)
	Include some of the server life circle
*/
package grpc

import (
	"crypto/md5"
	"fmt"
	"goim-pro/api/protos"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type server struct{}

func New() *grpc.Server {
	// TODO: should add the service discovery

	lis, err := net.Listen("tcp", "127.0.0.1:9090")
	if err != nil {
		log.Fatalf("TCP 监听失败: %v", err)
	} else {
		log.Printf("TCP 开始监听...", )
	}

	// 创建 gRPC 服务
	srv := grpc.NewServer()
	// 注册接口服务
	protos.RegisterWaiterServer(srv, &server{})
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

func (s *server) DoMD5(ctx context.Context, in *protos.Req) (*protos.Res, error) {
	fmt.Println("MD5方法请求的JSON: ", in.JsonStr)
	return &protos.Res{BackJson: "MD5 :" + fmt.Sprintf("%x", md5.Sum([]byte(in.JsonStr)))}, nil
}
