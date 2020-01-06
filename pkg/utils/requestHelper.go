package utils

import (
	"github.com/golang/protobuf/ptypes/any"
	"goim-pro/api/protos"
)

func NewReq (clientReq *protos.BaseClientRequest) map[string]interface{} {
	return nil
}
func NewResp(code int32, bytes []byte, message string) *protos.BaseServerResponse {
	return &protos.BaseServerResponse{
		Code: code,
		Data: &any.Any{
			TypeUrl: "",
			Value:   bytes,
		},
		Message: message,
	}
}
