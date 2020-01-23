package main

import (
	"crypto/md5"
	"fmt"
	"github.com/golang/protobuf/ptypes"
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

func (s *server) DoMD5(ctx context.Context, in *protos.Req) (*protos.Res, error) {
	fmt.Println("MD5方法请求的JSON: ", in)
	fmt.Println("Data: ", in.GetData())
	return &protos.Res{BackJson: "MD5 :" + fmt.Sprintf("%x", md5.Sum([]byte(in.JsonStr)))}, nil
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
	//protos.RegisterWaiterServer(s, &server{})

	s.RegisterService(&_Waiter_serviceDesc, &server{})

	// 在 gRPC 服务器上注册反射服务
	reflection.Register(s)

	// 将监听交给 gRPC 服务处理
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// TODO: for interceptor testing
var _Waiter_serviceDesc = grpc.ServiceDesc{
	ServiceName: "com.salty.protos.Waiter",
	HandlerType: (*protos.WaiterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DoMD5",
			Handler:    _Waiter_DoMD5_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "hellogrpc.proto",
}

func _Waiter_DoMD5_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(protos.GrpcReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/com.salty.protos.Waiter/DoMD5",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(protos.WaiterServer).DoMD5(ctx, req.(*protos.Req))
	}
	return interceptor(ctx, in, info, handler)
}
