package authsrv

import (
	"context"
	protos "goim-pro/api/protos/saltyv2"
	"goim-pro/internal/app/constants"
	"goim-pro/pkg/logs"
	"goim-pro/pkg/utils"
)

type smsServer struct{}

var logger = logs.GetLogger("INFO")

func New() protos.SMSServiceServer {
	return &smsServer{}
}

func (s *smsServer) ObtainSMSCode(ctx context.Context, req *protos.BasicReq) (res *protos.BasicResp, err error) {
	var smsReq protos.SMSReq
	err = utils.NewReq(req, &smsReq)
	if err != nil {
		logger.Errorf(`unmarshal error: %v`, err)
	} else {
		var code int32 = 200
		var verificationCode string = ""
		var msg string = "sending sms code success"
		switch smsReq.GetCodeType() {
		case protos.SMSReq_CodeType(constants.CodeTypeRegister):
			verificationCode = "123456"
			break
		case protos.SMSReq_CodeType(constants.CodeTypeLogin):
			verificationCode = "654321"
			break
		default:
			code = 400
			msg = "invalid request code type"
		}
		res = utils.NewResp(code, []byte(verificationCode), msg)
	}
	return
}
