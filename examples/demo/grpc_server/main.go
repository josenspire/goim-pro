package main

import (
	context2 "context"
	"crypto/md5"
	"fmt"
	example "goim-pro/api/protos/example"
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

func (s *server) DoMD5(ctx context.Context, in *example.Req) (*example.Res, error) {
	fmt.Println("MD5方法请求的JSON: ", in.JsonStr)
	fmt.Println("Data: ", in.GetData())
	return &example.Res{BackJson: "MD5 :" + fmt.Sprintf("%x", md5.Sum([]byte(in.JsonStr)))}, nil
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
	interceptor = func(ctx context2.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		reReq := req.(*example.Req)
		fmt.Println("=========", reReq.Data.Value)
		//req, _ = utils.NewMarshalAny(&example.Req{
		//	JsonStr:              "YYYY",
		//	Data:                 nil,
		//})
		reReq.Data.Value = []byte{1, 2, 3, 4, 5}
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
