package auth

import (
	"context"
	"fmt"
	protos "goim-pro/api/protos/salty"
	"goim-pro/internal/app/services/auth"
	errmsg "goim-pro/pkg/errors"
	"goim-pro/pkg/logs"
	"goim-pro/pkg/utils"
	"strings"
)

var (
	logger      = logs.GetLogger("INFO")
	authService authsrv.IAuthService
)

const (
	CodeLength = 6
)

type authServer struct {
}

func New() protos.SMSServiceServer {
	authService = authsrv.New()
	return &authServer{}
}

func (a *authServer) ObtainTelephoneSMSCode(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPC error) {
	resp, _ = utils.NewGRPCResp(protos.StatusCode_STATUS_OK, nil, "")

	var err error
	var smsReq protos.ObtainTelephoneSMSCodeReq
	if err = utils.UnMarshalAnyToMessage(req.GetData(), &smsReq); err != nil {
		logger.Errorf(`data unmarshal error: %s`, err.Error())
		resp.Code = protos.StatusCode_STATUS_BAD_REQUEST
		resp.Message = err.Error()
		return
	}

	if err = parameterCalibration(&smsReq); err != nil {
		resp.Code = protos.StatusCode_STATUS_BAD_REQUEST
		resp.Message = err.Error()
		return
	}
	telephone := smsReq.GetTelephone()
	operationType := smsReq.GetOperationType()

	code, tErr := authService.ObtainSMSCode(telephone, operationType)
	if tErr != nil {
		resp.Code = tErr.Code
		resp.Message = tErr.Detail
		return
	}
	resp.Message = fmt.Sprintf("send verification code succeed: %s", code)
	return
}

func (a *authServer) VerifyTelephoneSMSCode(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPC error) {
	resp, _ = utils.NewGRPCResp(protos.StatusCode_STATUS_OK, nil, "")

	var err error
	var verifyReq protos.VerifyTelephoneSMSCodeReq
	if err = utils.UnMarshalAnyToMessage(req.GetData(), &verifyReq); err != nil {
		logger.Errorf(`data unmarshal error: %s`, err.Error())
		resp.Code = protos.StatusCode_STATUS_BAD_REQUEST
		resp.Message = err.Error()
		return
	}

	operationType := verifyReq.GetOperationType()
	telephone := strings.Trim(verifyReq.GetTelephone(), "")
	codeStr := strings.Trim(verifyReq.SmsCode, "")

	if utils.IsContainEmptyString(telephone, codeStr) {
		resp.Code = protos.StatusCode_STATUS_BAD_REQUEST
		resp.Message = errmsg.ErrInvalidParameters.Error()
		return
	}

	_, tErr := authService.VerifySMSCode(telephone, operationType, codeStr)
	if tErr != nil {
		resp.Code = tErr.Code
		resp.Message = tErr.Detail
		return
	}

	obtainTelephoneSMSCodeResp := &protos.ObtainTelephoneSMSCodeResp{
		SmsCodeLength: CodeLength,
	}
	resp.Data, err = utils.MarshalMessageToAny(obtainTelephoneSMSCodeResp)
	if err != nil {
		logger.Errorf("obtain sms code response marshal message error: %s", err.Error())
	}
	resp.Message = "successful verification"
	return
}

func parameterCalibration(req *protos.ObtainTelephoneSMSCodeReq) (err error) {
	csErr := errmsg.ErrInvalidParameters
	req.Telephone = strings.Trim(req.GetTelephone(), "")

	if utils.IsEmptyStrings(req.GetTelephone()) {
		err = csErr
	}
	return
}
