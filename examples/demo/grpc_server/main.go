package main

import (
	"fmt"
	example "goim-pro/api/protos/example"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

const (
	address = "localhost:9090"
)

type server struct{}

func (server) SayHello(ctx context.Context, req *example.HelloReq) (resp *example.HelloResp, err error) {
	name := req.Name
	message := fmt.Sprintf("Hello %s, wellcome!!!", name)
	fmt.Println(message)

	resp = &example.HelloResp{
		Message: message,
	}
	return
}

func main() {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("监听失败: %v", err)
	} else {
		log.Printf("开始监听...")
	}

	var opts []grpc.ServerOption

	var interceptor grpc.UnaryServerInterceptor
	interceptor = func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		fmt.Println(info.FullMethod)
		return handler(ctx, req)
	}

	opts = append(opts, grpc.UnaryInterceptor(interceptor))

	s := grpc.NewServer(opts...) // 创建 gRPC 服务

	// 注册接口服务
	example.RegisterWaiterServer(s, &server{})

	// 在 gRPC 服务器上注册反射服务
	reflection.Register(s)

	// 将监听交给 gRPC 服务处理
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
