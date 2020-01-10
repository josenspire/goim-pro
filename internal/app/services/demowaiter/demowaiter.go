package demowaitersrv

import (
	"context"
	"crypto/md5"
	"fmt"
	example "goim-pro/api/protos/example"
	"goim-pro/pkg/logs"
)

type waiterServer struct{}

var logger = logs.GetLogger("INFO")

func New() example.WaiterServer {
	return &waiterServer{}
}

func (hw *waiterServer) DoMD5(ctx context.Context, in *example.Req) (*example.Res, error) {
	logger.Println("MD5方法请求的JSON: ", in.JsonStr)
	return &example.Res{BackJson: "MD5 :" + fmt.Sprintf("%x", md5.Sum([]byte(in.JsonStr)))}, nil
}
