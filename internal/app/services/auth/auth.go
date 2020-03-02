package authsrv

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v7"
	protos "goim-pro/api/protos/salty"
	. "goim-pro/internal/app/constants"
	"goim-pro/internal/app/repos"
	. "goim-pro/internal/app/repos/user"
	redsrv "goim-pro/pkg/db/redis"
	"goim-pro/pkg/logs"
	"goim-pro/pkg/utils"
	"net/http"
	"strings"
	"time"
)

type smsServer struct {
	userRepo IUserRepo
}

var (
	codeSize    = 6
	expiresTime = 15 * time.Minute // 15 min
	logger      = logs.GetLogger("INFO")
	myRedis     *redis.Client
)

func New() protos.SMSServiceServer {
	myRedis = redsrv.NewRedisConnection().GetRedisClient()

	repoServer := repos.New()
	return &smsServer{
		userRepo: repoServer.UserRepo,
	}
}

func (s *smsServer) ObtainSMSCode(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, grpcErr error) {
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
			resp.Code = http.StatusBadRequest
			resp.Message = "this telephone has been registered"
			return
		}

		redisKey = fmt.Sprintf("%d-%s", CodeTypeRegister, telephone)
		resp.Message = fmt.Sprintf("sending sms code success: %s", verificationCode)
	case protos.ObtainSMSCodeReq_CodeType(CodeTypeLogin):
		if !isTelephoneRegistered {
			resp.Code = http.StatusBadRequest
			resp.Message = "this telephone has not been registered"
			return
		}

		redisKey = fmt.Sprintf("%d-%s", CodeTypeLogin, telephone)
		resp.Message = fmt.Sprintf("sending sms code success: %s", verificationCode)
	case protos.ObtainSMSCodeReq_CodeType(CodeTypeResetPassword):
		if !isTelephoneRegistered {
			resp.Code = http.StatusBadRequest
			resp.Message = "this account does not exist"
			return
		}

		redisKey = fmt.Sprintf("%d-%s", CodeTypeResetPassword, telephone)
		resp.Message = fmt.Sprintf("sending sms code success: %s", verificationCode)
	default:
		resp.Code = http.StatusBadRequest
		resp.Message = "invalid request code type"
		return
	}
	if err = myRedis.Set(redisKey, verificationCode, expiresTime).Err(); err != nil {
		logger.Errorf("redis save error: %v", err)
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
	}
	return
}

func parameterCalibration(req *protos.ObtainSMSCodeReq) (err error) {
	csErr := errors.New("bad request, invalid parameters")
	req.Telephone = strings.Trim(req.GetTelephone(), "")

	if utils.IsEmptyStrings(req.GetTelephone()) {
		err = csErr
	}
	return
}
