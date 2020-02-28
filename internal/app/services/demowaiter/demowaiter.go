package demowaitersrv

import (
	"context"
	"crypto/md5"
	"fmt"
	"goim-pro/api/protos/example"
	"goim-pro/pkg/logs"
)

type waiterServer struct{}

var logger = logs.GetLogger("INFO")

func New() com_salty_protos.WaiterServer {
	return &waiterServer{}
}

func (hw *waiterServer) DoMD5(ctx context.Context, in *com_salty_protos.Req) (*com_salty_protos.Res, error) {
	logger.Println("MD5方法请求的JSON: ", in.JsonStr)
	return &com_salty_protos.Res{BackJson: "MD5 :" + fmt.Sprintf("%x", md5.Sum([]byte(in.JsonStr)))}, nil
}
