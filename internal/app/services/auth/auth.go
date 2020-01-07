package authsrv

import (
	"context"
	"encoding/json"
	"goim-pro/api/protos"
	"goim-pro/internal/app/constants"
	"goim-pro/pkg/logs"
	"goim-pro/pkg/utils"
)

type smsServer struct{}

var logger = logs.GetLogger("INFO")

func New() protos.SMSServiceServer {
	return &smsServer{}
}

func (s *smsServer) ObtainSMSCode(ctx context.Context, req *protos.BasicClientRequest) (res *protos.BasicServerResponse, err error) {
	reqBody := req.Data.GetValue()
	logger.Infoln(reqBody)

	var smsReq *protos.SMSReq
	err = json.Unmarshal(reqBody, &smsReq)
	if err != nil || smsReq == nil {
		logger.Errorf(`unmarshal error: %v`, err)
	} else {
		var code int32 = 200
		var verificationCode string = ""
		var msg string = "sending sms code success"
		switch smsReq.GetCodeType() {
		case protos.CodeType(constants.CodeTypeRegister):
			verificationCode = "123456"
			break
		case protos.CodeType(constants.CodeTypeLogin):
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
