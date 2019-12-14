package main

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

const (
	address = "localhost:9090"
)

func (s *server) DoMD5(ctx context.Context, in *protos.Req) (*protos.Res, error) {
	fmt.Println("MD5方法请求的JSON: ", in.JsonStr)
	return &protos.Res{BackJson: "MD5 :" + fmt.Sprintf("%x", md5.Sum([]byte(in.JsonStr)))}, nil
}

func main() {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("监听失败: %v", err)
	} else {
		log.Printf("开始监听...",)
	}

	s := grpc.NewServer() // 创建 gRPC 服务

	// 注册接口服务
	protos.RegisterWaiterServer(s, &server{})

	// 在 gRPC 服务器上注册反射服务
	reflection.Register(s)

	// 将监听交给 gRPC 服务处理
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
