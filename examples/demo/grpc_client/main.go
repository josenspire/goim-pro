package main

import (
	"context"
	protos "goim-pro/api/protos/example"
	"google.golang.org/grpc"
	"log"
)

const (
	address = "localhost:8080"
)

func main() {
	conn, err := grpc.Dial(
		address,
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("grpc connect fail: %v", err)
	}
	defer conn.Close()

	// create Writer service's client
	t := protos.NewWaiterClient(conn)

	//	调用 gRPC 接口
	tr, err := t.SayHello(context.Background(), &protos.GrpcReq{})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	//log.Printf("服务端响应：%s", tr.BackJson)
	log.Printf("服务端响应：%s", tr)
}
