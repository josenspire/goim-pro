package main

import (
	"context"
	protos "goim-pro/api/protos/salty"
	"goim-pro/pkg/utils"
	"google.golang.org/grpc"
	"log"
)

const (
	address = "111.231.238.209:9090"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("grpc connect fail: %v", err)
	}
	defer conn.Close()

	// create Writer service's client
	//t := protos.NewUserServiceClient(conn)
	t := protos.NewSMSServiceClient(conn)

	smsReq := protos.SMSReq{
		CodeType: protos.SMSReq_REGISTER,
		TargetAccount: &protos.SMSReq_Telephone{
			Telephone: "13631210000",
		},
	}
	anyData, _ := utils.MarshalMessageToAny(&smsReq)
	gprcReq := &protos.GrpcReq{
		Data: anyData,
	}

	// 调用 gRPC 接口
	tr, err := t.ObtainSMSCode(context.Background(), gprcReq)
	//tr, err := t.Register(context.Background(), gprcReq)
	if err != nil {
		log.Fatalf("could not greet: %v", err.Error())
	}
	log.Printf("服务端响应：%s", tr.GetData().Value)
}
