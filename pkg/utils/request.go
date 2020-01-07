package utils

import (
	"github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"goim-pro/api/protos"
	"goim-pro/pkg/logs"
)

var logger = logs.GetLogger("ERROR")

// unmarshal request: any
func NewReq(clientReq *protos.BasicClientRequest, pb proto.Message) (err error) {
	reqBody := clientReq.GetData()
	if err = ptypes.UnmarshalAny(reqBody, pb); err != nil {
		logger.Errorf("[ptypes] unmarshalAny error: %v", err)
		return
	}
	return err
}

// marshal response
func NewResp(code int32, bytes []byte, message string) *protos.BasicServerResponse {
	return &protos.BasicServerResponse{
		Code: code,
		Data: &any.Any{
			TypeUrl: "",
			Value:   bytes,
		},
		Message: message,
	}
}

func NewMarshalAny(pb proto.Message) (any *any.Any, err error) {
	any, err = ptypes.MarshalAny(pb)
	if err != nil {
		logger.Errorf("[ptypes] marshalAny error: %v", err)
	}
	return
}
