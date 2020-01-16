package authsrv

import (
	"context"
	"fmt"
	protos "goim-pro/api/protos/salty"
	"goim-pro/internal/app/constants"
	"goim-pro/pkg/logs"
	"goim-pro/pkg/utils"
	"net/http"
)

type smsServer struct{}

var logger = logs.GetLogger("INFO")

func New() protos.SMSServiceServer {
	return &smsServer{}
}

func (s *smsServer) ObtainSMSCode(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, err error) {
	resp, err = utils.NewGRPCResp(200, nil, "")

	var smsReq protos.SMSReq
	if err = utils.UnMarshalAnyToMessage(req.GetData(), &smsReq); err != nil {
		logger.Errorf(`data unmarshal error: %s`, err.Error())
		resp.Code = http.StatusBadRequest
		resp.Message = err.Error()
		return
	}

	resp.Message = "sending sms code success"
	var verificationCode string = ""
	switch smsReq.GetCodeType() {
	case protos.SMSReq_CodeType(constants.CodeTypeRegister):
		verificationCode = "123456"
		resp.Message = fmt.Sprintf("sending sms code success: %s", verificationCode)
		break
	case protos.SMSReq_CodeType(constants.CodeTypeLogin):
		verificationCode = "654321"
		resp.Message = fmt.Sprintf("sending sms code success: %s", verificationCode)
		break
	default:
		resp.Code = http.StatusBadRequest
		resp.Message = "invalid request code type"
	}
	return
}
