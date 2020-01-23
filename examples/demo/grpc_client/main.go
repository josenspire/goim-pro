package main

import (
	"context"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	protos "goim-pro/api/protos/salty"
	"goim-pro/pkg/utils"
	"google.golang.org/grpc"
	"log"
	"os"
)

const (
	address = "localhost:9090"
)

func main() {
	var clientInterceptor grpc.UnaryClientInterceptor
	clientInterceptor = func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		fmt.Println(req, method)

		var err error
		_req := req.(proto.Message)
		anyThing, _ := ptypes.MarshalAny(_req)

		err = invoker(ctx, method, &protos.GrpcReq{
			DeviceID: "asdfADF",
			Data:     anyThing,
		}, reply, cc, opts...)
		if err != nil {
			log.Println("接口调用出错", method, err)
			return err
		}
		return err
	}

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithUnaryInterceptor(clientInterceptor))
	if err != nil {
		log.Fatalf("grpc connect fail: %v", err)
	}
	defer conn.Close()

	// create Writer service's client
	t := protos.NewWaiterClient(conn)

	//	模拟请求数据
	res := "test123"
	//	os.Args[1] 为用户执行输入的参数 如：go run ***.go 123
	if len(os.Args) > 1 {
		res = os.Args[1]
	}

	any, _ := utils.MarshalMessageToAny(&protos.Req{
		JsonStr: "XXX",
	})
	//	调用 gRPC 接口
	tr, err := t.DoMD5(context.Background(), &protos.Req{
		JsonStr: res,
		Data:    any,
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	//log.Printf("服务端响应：%s", tr.BackJson)
	log.Printf("服务端响应：%s", tr)
}
