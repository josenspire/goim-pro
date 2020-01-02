package demowaitersrv

import (
	"context"
	"crypto/md5"
	"fmt"
	protos "goim-pro/api/protos"
	"goim-pro/pkg/logs"
)

type waiterServer struct {}

var logger = logs.GetLogger("INFO")

func New() protos.WaiterServer {
	return &waiterServer{}
}

func (hw *waiterServer) DoMD5(ctx context.Context, in *protos.Req) (*protos.Res, error) {
	logger.Println("MD5方法请求的JSON: ", in.JsonStr)
	return &protos.Res{BackJson: "MD5 :" + fmt.Sprintf("%x", md5.Sum([]byte(in.JsonStr)))}, nil
}
