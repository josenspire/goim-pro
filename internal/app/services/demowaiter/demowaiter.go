package demowaitersrv

import (
	"context"
	"fmt"
	example "goim-pro/api/protos/example"
	"goim-pro/pkg/http"
	"goim-pro/pkg/logs"
	"goim-pro/pkg/utils"
)

type waiterServer struct{}

var logger = logs.GetLogger("INFO")

func New() example.WaiterServer {
	return &waiterServer{}
}

func (waiterServer) SayHello(ctx context.Context, req *example.GrpcReq) (resp *example.GrpcResp, gRPCErr error) {
	resp = &example.GrpcResp{}

	var err error
	var demoReq example.HelloReq
	if err = utils.UnMarshalAnyToMessage(req.GetData(), &demoReq); err != nil {
		logger.Errorf(`data unmarshal error: %s`, err.Error())
		resp.Code = http.StatusBadRequest
		resp.Message = err.Error()
		return
	}

	name := demoReq.Name
	message := fmt.Sprintf("Hello %s, wellcome!!!", name)
	logger.Info(message)

	demoResp := &example.HelloResp{
		Message: message,
	}
	resp.Data, err = utils.MarshalMessageToAny(demoResp)
	if err != nil {
		logger.Errorf("register response marshal message error: %s", err.Error())
	}
	resp.Message = "user registration successful"

	return
}
