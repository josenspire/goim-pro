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
	"goim-pro/pkg/logs"
	"goim-pro/pkg/utils"
	"strings"
)

var (
	codeSize = 6
	logger   = logs.GetLogger("INFO")
	myRedis  redsrv.IMyRedis
	mysqlDB  *gorm.DB

	userRepo IUserRepo
)

type IAuthService interface {
	ObtainSMSCode(telephone string, operationType protos.SMSOperationType) (code string, tErr *TError)
	VerifySMSCode(telephone string, operationType protos.SMSOperationType, codeStr, deviceId string) (isPass bool, tErr *TError)
}

type AuthService struct {
}

func New() IAuthService {
	myRedis = redsrv.NewRedis()
	mysqlDB = mysqlsrv.NewMysql()
	userRepo = NewUserRepo(mysqlDB)
	return &AuthService{}
}

// ObtainSMSCode - obtain verification code by telephone with operation type
// operation type include:
// - register: SMSOperationType_REGISTER
// - login: SMSOperationType_LOGIN
// - reset password: SMSOperationType_RESET_PASSWORD
func (s *AuthService) ObtainSMSCode(telephone string, operationType protos.SMSOperationType) (code string, tErr *TError) {
	// checking user account status by telephone
	isTelephoneRegistered, err := userRepo.IsTelephoneOrEmailRegistered(telephone, "")
	if err != nil {
		logger.Errorf("checking sms telephone error: %s", err)
		return "", NewTError(protos.StatusCode_STATUS_INTERNAL_SERVER_ERROR, err)
	}

	// generate random num string
	verificationCode := utils.GenerateRandomNum(codeSize)
	var redisKey string = ""

	switch operationType {
	case protos.SMSOperationType_REGISTER:
		if isTelephoneRegistered {
			return "", NewTError(protos.StatusCode_STATUS_ACCOUNT_EXISTS, errmsg.ErrTelephoneExists)
		}

		verificationCode = "123401"
		redisKey = fmt.Sprintf("%d-%s", CodeTypeRegister, telephone)
		code = verificationCode
	case protos.SMSOperationType_LOGIN:
		if !isTelephoneRegistered {
			return "", NewTError(protos.StatusCode_STATUS_BAD_REQUEST, errmsg.ErrTelephoneNotExists)
		}

		verificationCode = "123402"
		redisKey = fmt.Sprintf("%d-%s", CodeTypeLogin, telephone)
		code = verificationCode
	case protos.SMSOperationType_RESET_PASSWORD:
		if !isTelephoneRegistered {
			return "", NewTError(protos.StatusCode_STATUS_ACCOUNT_NOT_EXISTS, errmsg.ErrAccountNotExists)
		}

		verificationCode = "123403"
		redisKey = fmt.Sprintf("%d-%s", CodeTypeResetPassword, telephone)
		code = verificationCode
	default:
		return "", NewTError(protos.StatusCode_STATUS_BAD_REQUEST, errmsg.ErrInvalidParameters)
	}
	if err = myRedis.RSet(redisKey, verificationCode, MinuteOf15); err != nil {
		logger.Errorf("redis save error: %v", err)
		return "", NewTError(protos.StatusCode_STATUS_INTERNAL_SERVER_ERROR, err)
	}
	return code, nil
}

// VerifySMSCode - verify sms code
func (s *AuthService) VerifySMSCode(telephone string, operationType protos.SMSOperationType, codeStr, deviceId string) (isPass bool, tErr *TError) {
	codeKey := fmt.Sprintf("%d-%s", operationType, telephone)

	code := myRedis.RGet(codeKey)
	if utils.IsEmptyStrings(code) {
		return false, NewTError(protos.StatusCode_STATUS_BAD_REQUEST, errmsg.ErrVerificationCodeExpired)
	}

	if !strings.EqualFold(codeStr, code) {
		return false, NewTError(protos.StatusCode_STATUS_BAD_REQUEST, errmsg.ErrInvalidVerificationCode)
	}

	if operationType == protos.SMSOperationType_LOGIN {
		if deviceId == "" {
			return false, NewTError(protos.StatusCode_STATUS_ABNORMAL_DEVICE_INFO, errmsg.ErrAbnormalDeviceInfo)
		}
		err := refreshUserProfileWithDeviceId(telephone, deviceId)
		if err == gorm.ErrRecordNotFound {
			return false, NewTError(protos.StatusCode_STATUS_BAD_REQUEST, errmsg.ErrTelephoneNotExists)
		} else if err != nil {
			return false, NewTError(protos.StatusCode_STATUS_INTERNAL_SERVER_ERROR, errmsg.ErrSystemUncheckException)
		}
	}

	myRedis.RDel(codeKey)
	return true, nil
}

func refreshUserProfileWithDeviceId(telephone, deviceId string) (err error) {
	var userCriteria = map[string]interface{}{
		"Telephone": telephone,
	}
	var updated = map[string]interface{}{
		"DeviceId": deviceId,
	}
	if err := userRepo.FindOneAndUpdateProfile(userCriteria, updated); err != nil {
		return err
	}
	return nil
}
