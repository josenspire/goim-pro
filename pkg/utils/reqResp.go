package utils

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	protos "goim-pro/api/protos/salty"
)

// unmarshal request: any
func UnmarshalGRPCReq(clientReq *protos.GrpcReq, pb proto.Message) (err error) {
	anyReqBody := clientReq.GetData()
	if err = UnMarshalAnyToMessage(anyReqBody, pb); err != nil {
		logger.Errorf("unmarshal gRPC request error: %s", err.Error())
	}
	return
}

// marshal response
func MarshalGRPCRespWithBytes(code int32, bytes []byte, message string) *protos.GrpcResp {
	return &protos.GrpcResp{
		Code: code,
		Data: &any.Any{
			TypeUrl: "",
			Value:   bytes,
		},
		Message: message,
	}
}

// marshal response
func NewGRPCResp(code int32, pb proto.Message, message string) (gRPCResp *protos.GrpcResp, err error) {
	var anyData *any.Any
	if pb != nil {
		if anyData, err = MarshalMessageToAny(pb); err != nil {
			logger.Errorf("[NewGRPCResp] marshal message error: %s", err.Error())
			return
		}
	}
	gRPCResp = &protos.GrpcResp{
		Code:    code,
		Data:    anyData,
		Message: message,
	}
	return
}
