package utils

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	protos "goim-pro/api/protos/saltyv2"
	"goim-pro/pkg/logs"
)

var logger = logs.GetLogger("ERROR")

// unmarshal request: any
func NewReq(clientReq *protos.BasicReq, pb proto.Message) (err error) {
	reqBody := clientReq.GetData()
	if err = ptypes.UnmarshalAny(reqBody, pb); err != nil {
		logger.Errorf("[ptypes] unmarshalAny error: %v", err)
		return
	}
	return err
}

// marshal response
func NewResp(code int32, bytes []byte, message string) *protos.BasicResp {
	return &protos.BasicResp{
		Code: code,
		Data: &any.Any{
			TypeUrl: "",
			Value:   bytes,
		},
		Message: message,
	}
}
