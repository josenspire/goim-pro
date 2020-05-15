package authsrv

import (
	"fmt"
	"github.com/jinzhu/gorm"
	protos "goim-pro/api/protos/salty"
	. "goim-pro/internal/app/constants"
	. "goim-pro/internal/app/models/errors"
	. "goim-pro/internal/app/repos/user"
	mysqlsrv "goim-pro/pkg/db/mysql"
	redsrv "goim-pro/pkg/db/redis"
	"goim-pro/pkg/errors"
	"goim-pro/pkg/http"
	"goim-pro/pkg/logs"
	"goim-pro/pkg/utils"
)

var (
	codeSize = 6
	logger   = logs.GetLogger("INFO")
	myRedis  *redsrv.BaseClient
	mysqlDB  *gorm.DB

	userRepo IUserRepo
)

type AuthService struct {
}

func New() *AuthService {
	myRedis = redsrv.NewRedisConnection().GetRedisClient()
	mysqlDB = mysqlsrv.NewMysqlConnection().GetMysqlInstance()
	userRepo = NewUserRepo(mysqlDB)
	return &AuthService{}
}

func (s *AuthService) ObtainSMSCode(telephone string, codeType protos.ObtainSMSCodeReq_CodeType) (code string, tErr *TError) {
	// checking user account status by telephone
	isTelephoneRegistered, err := userRepo.IsTelephoneOrEmailRegistered(telephone, "")
	if err != nil {
		logger.Errorf("checking sms telephone error: %s", err)
		return "", NewTError(http.StatusInternalServerError, err)
	}

	// generate random num string
	verificationCode := utils.GenerateRandomNum(codeSize)
	var redisKey string = ""

	switch codeType {
	case protos.ObtainSMSCodeReq_CodeType(CodeTypeRegister):
		if isTelephoneRegistered {
			return "", NewTError(http.StatusAccountExists, errmsg.ErrTelephoneExists)
		}

		verificationCode = "123401"
		redisKey = fmt.Sprintf("%d-%s", CodeTypeRegister, telephone)
		code = verificationCode
	case protos.ObtainSMSCodeReq_CodeType(CodeTypeLogin):
		if !isTelephoneRegistered {
			return "", NewTError(http.StatusBadRequest, errmsg.ErrTelephoneNotExists)
		}

		verificationCode = "123402"
		redisKey = fmt.Sprintf("%d-%s", CodeTypeLogin, telephone)
		code = verificationCode
	case protos.ObtainSMSCodeReq_CodeType(CodeTypeResetPassword):
		if !isTelephoneRegistered {
			return "", NewTError(http.StatusBadRequest, errmsg.ErrAccountNotExists)
		}

		verificationCode = "123403"
		redisKey = fmt.Sprintf("%d-%s", CodeTypeResetPassword, telephone)
		code = verificationCode
	default:
		return "", NewTError(http.StatusBadRequest, errmsg.ErrInvalidParameters)
	}
	if err = myRedis.Set(redisKey, verificationCode, MinuteOf15); err != nil {
		logger.Errorf("redis save error: %v", err)
		return "", NewTError(http.StatusInternalServerError, err)
	}
	return code, nil
}
