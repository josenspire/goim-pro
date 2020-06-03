package usersrv

import (
	"fmt"
	"github.com/jinzhu/gorm"
	protos "goim-pro/api/protos/salty"
	"goim-pro/config"
	. "goim-pro/internal/app/constants"
	"goim-pro/internal/app/models"
	. "goim-pro/internal/app/models/errors"
	. "goim-pro/internal/app/repos/user"
	"goim-pro/internal/app/services/converters"
	mysqlsrv "goim-pro/pkg/db/mysql"
	redsrv "goim-pro/pkg/db/redis"
	"goim-pro/pkg/errors"
	"goim-pro/pkg/logs"
	"goim-pro/pkg/utils"
)

var (
	logger = logs.GetLogger("INFO")

	myRedis redsrv.IMyRedis
	mysqlDB *gorm.DB

	userRepo IUserRepo
)

type UserService struct {
}

func New() *UserService {
	myRedis = redsrv.NewRedis()
	mysqlDB = mysqlsrv.NewMysql()
	userRepo = NewUserRepo(mysqlDB)

	return &UserService{}
}

func (s *UserService) Register(userProfile *protos.UserProfile, password string) (tErr *TError) {
	var telephone = userProfile.Telephone
	var email = userProfile.Email

	userProfile.UserId = utils.NewULID()

	isRegistered, err := userRepo.IsTelephoneOrEmailRegistered(telephone, email)
	if err != nil {
		logger.Errorf("checking telephone validity error: %v", err.Error())
		return NewTError(protos.StatusCode_STATUS_INTERNAL_SERVER_ERROR, err)
	}
	if isRegistered {
		return NewTError(protos.StatusCode_STATUS_ACCOUNT_EXISTS, errmsg.ErrAccountAlreadyExists)
	}
	if err = userRepo.Register(&models.User{
		Password:    password,
		UserProfile: converters.ConvertProto2EntityForUserProfile(userProfile),
	}); err != nil {
		logger.Errorf("register user error: %v", err.Error())
		return NewTError(protos.StatusCode_STATUS_INTERNAL_SERVER_ERROR, err)
	}

	return nil
}

func (s *UserService) Login(telephone, email, enPassword, deviceId string, osVersion protos.GrpcReq_OS) (user *models.User, token string, tErr *TError) {
	var isNeedVerify bool = false
	var err error
	isNeedVerify, user, err = accountLogin(telephone, email, enPassword, deviceId, osVersion)
	if err != nil {
		return nil, "", NewTError(protos.StatusCode_STATUS_BAD_REQUEST, err)
	}
	if isNeedVerify {
		return nil, "", NewTError(protos.StatusCode_STATUS_ACCOUNT_AUTHORIZED_REQUIRED, errmsg.ErrAccountSecurityVerification)
	}
	// gen and save token
	token = utils.NewToken([]byte(user.UserId))
	if err = myRedis.RSet(fmt.Sprintf("TK-%s", user.UserId), token, ThreeDays); err != nil {
		logger.Errorf("redis save token error: %v", err)
		return nil, "", NewTError(protos.StatusCode_STATUS_INTERNAL_SERVER_ERROR, err)
	}
	return user, token, nil
}

func (s *UserService) Logout(token string, isMandatoryLogout bool) (tErr *TError) {
	isValid, payload, err := utils.TokenVerify(token)
	if err != nil {
		logger.Errorf("logout by token error: %s", err.Error())
		return NewTError(protos.StatusCode_STATUS_INTERNAL_SERVER_ERROR, err)
	}
	if !isValid {
		return NewTError(protos.StatusCode_STATUS_BAD_REQUEST, errmsg.ErrRepeatOperation)
	}

	tokenKey := fmt.Sprintf("TK-%s", string(payload))
	_token := myRedis.RGet(tokenKey)
	if _token == "" {
		return NewTError(protos.StatusCode_STATUS_BAD_REQUEST, errmsg.ErrRepeatOperation)
	}
	// TODO: if true, mandatory and will remove all online user
	if isMandatoryLogout {
		myRedis.RDel(tokenKey)
	} else {
		myRedis.RDel(tokenKey)
	}
	return nil
}

func (s *UserService) UpdateUserInfo(userId string, userProfile *models.UserProfile) (tErr *TError) {
	originUserProfile, err := userRepo.FindByUserId(userId)
	if err != nil {
		return NewTError(protos.StatusCode_STATUS_BAD_REQUEST, errmsg.ErrInvalidParameters)
	}
	if originUserProfile == nil {
		return NewTError(protos.StatusCode_STATUS_BAD_REQUEST, errmsg.ErrInvalidUserId)
	}
	// nothing change, don't need to update
	if isProfileNothing2Update(&originUserProfile.UserProfile, userProfile) {
		return
	}

	criteria := map[string]interface{}{
		"userId": userId,
	}
	//updateMap := utils.TransformStructToMap(userProfile)
	//utils.RemoveMapProperties(updateMap, "UserId", "Telephone", "Email", "Avatar")
	updated := map[string]interface{}{
		"description": userProfile.Description,
		"location":    userProfile.Location,
		"birthday":    userProfile.Birthday,
		"sex":         userProfile.Sex,
		"nickname":    userProfile.Nickname,
	}
	err = userRepo.FindOneAndUpdateProfile(criteria, updated)
	if err != nil {
		logger.Errorf("[get user info] response marshal message error: %s", err.Error())
		return NewTError(protos.StatusCode_STATUS_INTERNAL_SERVER_ERROR, err)
	}
	return
}

// ResetPassword - reset user's password by verification code
func (s *UserService) ResetPassword(telephone, email, newPassword string) (tErr *TError) {
	isRegistered, err := userRepo.IsTelephoneOrEmailRegistered(telephone, email)
	if err != nil {
		logger.Errorf("reset password error: %s", err.Error())
		return NewTError(protos.StatusCode_STATUS_BAD_REQUEST, err)
	}
	if !isRegistered {
		logger.Warn("reset password failed: account not exist")
		return NewTError(protos.StatusCode_STATUS_BAD_REQUEST, errmsg.ErrUserNotExists)
	}

	// reset password by old password
	//if oldPassword == newPassword {
	//	return NewTError(protos.StatusCode_STATUS_BAD_REQUEST, errmsg.ErrRepeatPassword)
	//}
	//var user *models.User
	//enPassword := utils.NewSHA256(oldPassword, config.GetApiSecretKey())
	//if telephone != "" {
	//	user, err = userRepo.QueryByTelephoneAndPassword(telephone, enPassword)
	//} else {
	//	user, err = userRepo.QueryByEmailAndPassword(email, enPassword)
	//}
	//if err != nil {
	//	logger.Errorf("reset password error: %s", err.Error())
	//	return NewTError(protos.StatusCode_STATUS_BAD_REQUEST, err)
	//}
	//if user == nil {
	//	return NewTError(protos.StatusCode_STATUS_BAD_REQUEST, errmsg.ErrAccountOrPwdInvalid)
	//}

	//if verificationCode != "" {
	//	// 1. code + newPassword
	//	code := myRedis.RGet(fmt.Sprintf("%d-%s", CodeTypeResetPassword, telephone))
	//	if verificationCode != string(code) {
	//		return NewTError(protos.StatusCode_STATUS_BAD_REQUEST, errmsg.ErrInvalidVerificationCode)
	//	}
	//} else {
	//	// 2. oldPassword + newPassword
	//	if oldPassword == newPassword {
	//		return NewTError(protos.StatusCode_STATUS_BAD_REQUEST, errmsg.ErrRepeatPassword)
	//	}
	//	var user *models.User
	//	enPassword := utils.NewSHA256(oldPassword, config.GetApiSecretKey())
	//	if telephone != "" {
	//		user, err = userRepo.QueryByTelephoneAndPassword(telephone, enPassword)
	//	} else {
	//		user, err = userRepo.QueryByEmailAndPassword(email, enPassword)
	//	}
	//	if err != nil {
	//		logger.Errorf("reset password error: %s", err.Error())
	//		return NewTError(protos.StatusCode_STATUS_BAD_REQUEST, err)
	//	}
	//	if user == nil {
	//		return NewTError(protos.StatusCode_STATUS_BAD_REQUEST, errmsg.ErrAccountOrPwdInvalid)
	//	}
	//}

	enNewPassword := utils.NewSHA256(newPassword, config.GetApiSecretKey())
	if telephone != "" {
		if err = userRepo.ResetPasswordByTelephone(telephone, enNewPassword); err != nil {
			logger.Errorf("reset password error: %s", err.Error())
			return NewTError(protos.StatusCode_STATUS_INTERNAL_SERVER_ERROR, err)
		}
		//myRedis.RDel(fmt.Sprintf("%d-%s", CodeTypeResetPassword, telephone))
	} else {
		if err = userRepo.ResetPasswordByEmail(email, enNewPassword); err != nil {
			logger.Errorf("reset password error: %s", err.Error())
			return NewTError(protos.StatusCode_STATUS_INTERNAL_SERVER_ERROR, err)
		}
	}

	return nil
}

func (s *UserService) GetUserInfo(userId string) (profile *models.UserProfile, tErr *TError) {
	user, err := userRepo.FindByUserId(userId)
	if err != nil {
		return nil, NewTError(protos.StatusCode_STATUS_INTERNAL_SERVER_ERROR, err)
	}
	if user == nil {
		return nil, NewTError(protos.StatusCode_STATUS_BAD_REQUEST, errmsg.ErrInvalidUserId)
	}
	return &user.UserProfile, nil
}

// QueryUserInfo - query user information by telephone, email
func (s *UserService) QueryUserInfo(telephone, email string) (profile *models.UserProfile, tErr *TError) {
	userCriteria := make(map[string]interface{})
	if utils.IsEmptyStrings(telephone) {
		userCriteria["email"] = email
	} else {
		userCriteria["telephone"] = telephone
	}

	user, err := userRepo.FindOneUser(userCriteria)
	if err != nil {
		return nil, NewTError(protos.StatusCode_STATUS_INTERNAL_SERVER_ERROR, err)
	}
	if user == nil {
		return nil, NewTError(protos.StatusCode_STATUS_BAD_REQUEST, errmsg.ErrUserNotExists)
	}

	return &user.UserProfile, nil
}

func accountLogin(telephone, email, enPassword, deviceId string, osVersion protos.GrpcReq_OS) (isNeedVerify bool, user *models.User, err error) {
	var isLoginByTelephone bool = false
	if telephone != "" {
		isLoginByTelephone = true
	}

	user, err = loginWithPassword(isLoginByTelephone, telephone, email, enPassword)
	if err != nil {
		return false, nil, err
	}
	if user == nil {
		return false, nil, errmsg.ErrAccountOrPwdInvalid
	}
	// TODO: need to had previous verify function
	//return isNeedToSMSVerify(deviceId, osVersion, user), user, nil
	return false, user, nil
}

func isProfileNothing2Update(originProfile, newProfile *models.UserProfile) bool {
	return utils.DeepEqual(originProfile, newProfile)
}

func loginWithPassword(isTelephone bool, telephone, email, enPassword string) (user *models.User, err error) {
	if isTelephone {
		return userRepo.QueryByTelephoneAndPassword(telephone, enPassword)
	} else {
		return userRepo.QueryByEmailAndPassword(email, enPassword)
	}
}

func loginWithVerificationCode(isTelephone bool, telephone, email, verificationCode string) (user *models.User, err error) {
	var codeKey string = ""
	var criteria = make(map[string]interface{})
	if isTelephone {
		codeKey = fmt.Sprintf("%d-%s", CodeTypeLogin, telephone)
		criteria = map[string]interface{}{
			"telephone": telephone,
		}
	} else {
		codeKey = fmt.Sprintf("%d-%s", CodeTypeLogin, email)
		criteria = map[string]interface{}{
			"email": email,
		}
	}
	code := myRedis.RGet(codeKey)
	if verificationCode == string(code) {
		if user, err = userRepo.FindOneUser(criteria); err != nil {
			logger.Errorf("find user by account error: %s", err.Error())
			return nil, err
		}
		myRedis.RDel(codeKey)
		return user, nil
	}
	return nil, errmsg.ErrInvalidVerificationCode
}

func isNeedToSMSVerify(deviceId string, osVersion protos.GrpcReq_OS, orgUser *models.User) bool {
	// TODO: should divide strategy by platform
	if deviceId == "" || deviceId != orgUser.DeviceId {
		return true
	}

	return false

	//switch osVersion {
	//case protos.GrpcReq_ANDROID:
	//case protos.GrpcReq_IOS:
	//case protos.GrpcReq_WINDOWS:
	//case protos.GrpcReq_UNKNOWN:
	//}
}
