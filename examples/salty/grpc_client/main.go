package main

import (
	"context"
	"goim-pro/api/protos"
	"google.golang.org/grpc"
	"log"
	"os"
)

const (
	address = "localhost:9090"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
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

	//	调用 gRPC 接口
	tr, err := t.DoMD5(context.Background(), &protos.Req{JsonStr: res})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("服务端响应：%s", tr.BackJson)
}
