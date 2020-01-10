package main

import (
	"context"
	"encoding/json"
	"github.com/golang/protobuf/ptypes/any"
	protos "goim-pro/api/protos/saltyv2"
	"google.golang.org/grpc"
	"log"
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
	t := protos.NewUserServiceClient(conn)

	userData := map[string]interface{}{
		"username": "JAMES",
		"password": "1234567890",
		"age":      24,
	}
	dataByte, _ := json.Marshal(userData)
	anyData := &any.Any{
		Value: dataByte,
	}

	reqBody := &protos.BasicReq{
		Data: anyData,
	}

	// 调用 gRPC 接口
	tr, err := t.Register(context.Background(), reqBody)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("服务端响应：%s", tr.GetData().Value)
}
