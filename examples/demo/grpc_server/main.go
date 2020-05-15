package main

import (
	"crypto/md5"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	demo "goim-pro/api/protos/example"
	protos "goim-pro/api/protos/salty"
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
		fmt.Println("----------", req)
		//var unmarshaler proto.Unmarshaler
		reReq := req.(*protos.GrpcReq)
		log.Println(reReq)

		var pb protos.Req
		if err = ptypes.UnmarshalAny(reReq.GetData(), &pb); err != nil {
			//if err = ptypes.UnmarshalAny(any, pb); err != nil {
			fmt.Println(err)
		}
		return handler(ctx, &pb)
	}

	opts = append(opts, grpc.UnaryInterceptor(interceptor))

	s := grpc.NewServer(opts...) // 创建 gRPC 服务

	// 注册接口服务
	demo.RegisterWaiterServer(s, &server{})

	// 在 gRPC 服务器上注册反射服务
	reflection.Register(s)

	// 将监听交给 gRPC 服务处理
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
