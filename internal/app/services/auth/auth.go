package authsrv

import (
	"context"
	"fmt"
	protos "goim-pro/api/protos/salty"
	"goim-pro/internal/app/constants"
	"goim-pro/pkg/logs"
	"goim-pro/pkg/utils"
)

type smsServer struct{}

var logger = logs.GetLogger("INFO")

func New() protos.SMSServiceServer {
	return &smsServer{}
}

func (s *smsServer) ObtainSMSCode(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, err error) {
	resp = utils.NewResp(200, nil, "")

	var smsReq protos.SMSReq
	err = utils.NewReq(req, &smsReq)
	if err != nil {
		resp.Code = 500
		resp.Message = err.Error()
		logger.Errorf(`unmarshal error: %v`, err)
		return
	} else {
		resp.Code = 200
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
			resp.Code = 400
			resp.Message = "invalid request code type"
		}
		resp.Data, _ = utils.MarshalMessageToAny(&protos.SMSResp{})
	}
	return
}
