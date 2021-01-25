package demowaitersrv

import (
	"context"
	"fmt"
	example "goim-pro/api/protos/example"
	"goim-pro/pkg/logs"
)

type waiterServer struct{}

var logger = logs.GetLogger("INFO")

func New() example.WaiterServer {
	return &waiterServer{}
}

func (waiterServer) SayHello(ctx context.Context, req *example.HelloReq) (resp *example.HelloResp, gRPCErr error) {

	name := req.Name
	message := fmt.Sprintf("Hello %s, wellcome!!!", name)
	logger.Info(message)

	resp = &example.HelloResp{
		Message: message,
	}

	return
}
