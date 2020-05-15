package auth

import (
	"context"
	"fmt"
	protos "goim-pro/api/protos/salty"
	"goim-pro/internal/app/services/auth"
	errmsg "goim-pro/pkg/errors"
	"goim-pro/pkg/http"
	"goim-pro/pkg/logs"
	"goim-pro/pkg/utils"
	"strings"
)

var (
	logger      = logs.GetLogger("INFO")
	authService *authsrv.AuthService
)

type authServer struct {
}

func New() protos.SMSServiceServer {
	authService = authsrv.New()
	return &authServer{}
}

func (a authServer) ObtainSMSCode(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	resp, _ = utils.NewGRPCResp(http.StatusOK, nil, "")

	var err error
	var smsReq protos.ObtainSMSCodeReq
	if err = utils.UnMarshalAnyToMessage(req.GetData(), &smsReq); err != nil {
		logger.Errorf(`data unmarshal error: %s`, err.Error())
		resp.Code = http.StatusBadRequest
		resp.Message = err.Error()
		return
	}

	if err = parameterCalibration(&smsReq); err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = err.Error()
		return
	}
	telephone := smsReq.GetTelephone()
	codeType := smsReq.GetCodeType()

	code, tErr := authService.ObtainSMSCode(telephone, codeType)
	if tErr != nil {
		resp.Code = tErr.Code
		resp.Message = tErr.Detail
		return
	}
	resp.Message = fmt.Sprintf("send verification code succeed: %s", code)
	return
}

func parameterCalibration(req *protos.ObtainSMSCodeReq) (err error) {
	csErr := errmsg.ErrInvalidParameters
	req.Telephone = strings.Trim(req.GetTelephone(), "")

	if utils.IsEmptyStrings(req.GetTelephone()) {
		err = csErr
	}
	return
}
