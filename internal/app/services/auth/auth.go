package authsrv

import (
	"context"
	"fmt"
	"github.com/jinzhu/gorm"
	protos "goim-pro/api/protos/salty"
	. "goim-pro/internal/app/constants"
	. "goim-pro/internal/app/repos/user"
	mysqlsrv "goim-pro/pkg/db/mysql"
	redsrv "goim-pro/pkg/db/redis"
	"goim-pro/pkg/http"
	"goim-pro/pkg/logs"
	"goim-pro/pkg/utils"
	"strings"
)

var (
	codeSize = 6
	logger   = logs.GetLogger("INFO")
	myRedis  *redsrv.BaseClient
	mysqlDB  *gorm.DB
)

type smsService struct {
	userRepo IUserRepo
}

func New() protos.SMSServiceServer {
	myRedis = redsrv.NewRedisConnection().GetRedisClient()
	mysqlDB = mysqlsrv.NewMysqlConnection().GetMysqlInstance()

	return &smsService{
		userRepo: NewUserRepo(mysqlDB),
	}
}

func (s *smsService) ObtainSMSCode(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, grpcErr error) {
	resp, _ = utils.NewGRPCResp(http.StatusOK, nil, "")

	var smsReq protos.ObtainSMSCodeReq
	var err error
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

	// checking user account status by telephone
	isTelephoneRegistered, err := s.userRepo.IsTelephoneOrEmailRegistered(telephone, "")
	if err != nil {
		logger.Errorf("checking sms telephone error: %s", err)
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		return
	}

	// generate random num string
	verificationCode := utils.GenerateRandomNum(codeSize)
	var redisKey string = ""

	resp.Message = "sending sms code success"
	switch codeType {
	case protos.ObtainSMSCodeReq_CodeType(CodeTypeRegister):
		if isTelephoneRegistered {
			resp.Code = http.StatusAccountExists
			resp.Message = "this telephone has been registered"
			return
		}

		verificationCode = "123401"
		redisKey = fmt.Sprintf("%d-%s", CodeTypeRegister, telephone)
		resp.Message = fmt.Sprintf("sending sms code success: %s", verificationCode)
	case protos.ObtainSMSCodeReq_CodeType(CodeTypeLogin):
		if !isTelephoneRegistered {
			resp.Code = http.StatusBadRequest
			resp.Message = "this telephone has not been registered"
			return
		}

		verificationCode = "123402"
		redisKey = fmt.Sprintf("%d-%s", CodeTypeLogin, telephone)
		resp.Message = fmt.Sprintf("sending sms code success: %s", verificationCode)
	case protos.ObtainSMSCodeReq_CodeType(CodeTypeResetPassword):
		if !isTelephoneRegistered {
			resp.Code = http.StatusBadRequest
			resp.Message = "this account does not exist"
			return
		}

		verificationCode = "123403"
		redisKey = fmt.Sprintf("%d-%s", CodeTypeResetPassword, telephone)
		resp.Message = fmt.Sprintf("sending sms code success: %s", verificationCode)
	default:
		resp.Code = http.StatusBadRequest
		resp.Message = "invalid request code type"
		return
	}
	if err = myRedis.Set(redisKey, verificationCode, MinuteOf15); err != nil {
		logger.Errorf("redis save error: %v", err)
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
	}
	return
}

func parameterCalibration(req *protos.ObtainSMSCodeReq) (err error) {
	csErr := utils.ErrInvalidParameters
	req.Telephone = strings.Trim(req.GetTelephone(), "")

	if utils.IsEmptyStrings(req.GetTelephone()) {
		err = csErr
	}
	return
}
