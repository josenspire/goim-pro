package usersrv

import (
	"context"
	"encoding/json"
	any "github.com/golang/protobuf/ptypes/any"
	"goim-pro/api/protos"
	"goim-pro/pkg/logs"
)

type userServer struct{}

var logger = logs.GetLogger("INFO")

func New() protos.UserServiceServer {
	return &userServer{}
}

func (us *userServer) Register(ctx context.Context, req *protos.BaseClientRequest) (res *protos.BaseServerResponse, err error) {
	reqBody := req.Data.GetValue()
	logger.Infoln(reqBody)

	reqMap := make(map[string]interface{})
	err = json.Unmarshal(reqBody, &reqMap)
	if err != nil {
		logger.Errorf(`unmarshal error: %v`, err)
	} else {
		if reqMap["username"] == "JAMES" && reqMap["password"] == "1234567890" {
			var msg = "user regist successful.."
			res = &protos.BaseServerResponse{
				Code: 200,
				Data: &any.Any{Value: []byte(msg)},
				Message: "",
			}
		} else {
			var msg = "username or password incorrect"
			res = &protos.BaseServerResponse{
				Code: 400,
				Data: &any.Any{Value: []byte(msg)},
				Message: msg,
			}
		}
	}
	return
}

func (us *userServer) Login(ctx context.Context, req *protos.BaseClientRequest) (*protos.BaseServerResponse, error) {
	panic("implement me")
}
