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

func (s *smsServer) ObtainSMSCode(ctx context.Context, req *protos.BaseClientRequest) (res *protos.BaseServerResponse, err error) {
	reqBody := req.Data.GetValue()
	logger.Infoln(reqBody)

	var smsReq *protos.SMSReq
	err = json.Unmarshal(reqBody, &smsReq)
	if err != nil || smsReq == nil {
		logger.Errorf(`unmarshal error: %v`, err)
	} else {
		var code int32 = 200
		var verifyCode string = ""
		var msg string = "sending sms code success"
		switch smsReq.GetCodeType() {
		case protos.CodeType(constants.CodeTypeRegister):
			verifyCode = "123456"
			break
		case protos.CodeType(constants.CodeTypeLogin):
			verifyCode = "654321"
			break
		default:
			code = 400
			msg = "invalid request code type"
		}
		res = utils.NewResp(code, []byte(verifyCode), msg)
	}
	return
}
