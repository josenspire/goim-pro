package demo

import (
	"context"
	example "goim-pro/api/protos/example"
	"goim-pro/pkg/logs"
	"goim-pro/pkg/utils"
)

var logger = logs.GetLogger("INFO")

func SayHello(t example.WaiterClient) {
	reqData := &example.HelloReq{
		Name: "JAMES",
	}
	anyData, _ := utils.MarshalMessageToAny(reqData)
	resp, err := t.SayHello(context.Background(), &example.GrpcReq{
		Data:     anyData,
	})
	if err != nil {
		logger.Error(err)
	}

	printResp(resp)
}

func printResp(resp *example.GrpcResp) {
	logger.Infof("[code]: %d", resp.GetCode())
	logger.Infof("[data]: %s", resp.GetData().GetValue())
	logger.Infof("[message]: %s", resp.GetMessage())
}
