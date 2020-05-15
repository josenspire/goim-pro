package usersrv

import (
	"context"
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
	"goim-pro/pkg/http"
	"goim-pro/pkg/logs"
	"goim-pro/pkg/utils"
	"strings"
)

var (
	logger = logs.GetLogger("INFO")

	myRedis *redsrv.BaseClient
	mysqlDB *gorm.DB

	userRepo IUserRepo
)

type UserService struct {
}

func New() *UserService {
	myRedis = redsrv.NewRedisConnection().GetRedisClient()
	mysqlDB = mysqlsrv.NewMysqlConnection().GetMysqlInstance()
	userRepo = NewUserRepo(mysqlDB)

	return &UserService{}
}

func (s *UserService) Register(userProfile *protos.UserProfile, password, verificationCode string) (tErr *TError) {
	var telephone = userProfile.Telephone
	var email = userProfile.Email

	isValid, err := isVerificationCodeValid(verificationCode, telephone)
	if !isValid {
		logger.Warnf("verification code is invalid: %s", verificationCode)
		return NewTError(http.StatusBadRequest, errmsg.ErrInvalidVerificationCode)
	}
	userProfile.UserId = utils.NewULID()

	isRegistered, err := userRepo.IsTelephoneOrEmailRegistered(telephone, email)
	if err != nil {
		logger.Errorf("checking telephone validity error: %v", err.Error())
		return NewTError(http.StatusInternalServerError, err)
	}
	if isRegistered {
		return NewTError(http.StatusAccountExists, errmsg.ErrAccountAlreadyExists)
	}
	if err = userRepo.Register(&models.User{
		Password:    password,
		UserProfile: converters.ConvertProto2EntityForUserProfile(userProfile),
	}); err != nil {
		logger.Errorf("register user error: %v", err.Error())
		return NewTError(http.StatusInternalServerError, err)
	}
	// remove cache
	myRedis.Del(fmt.Sprintf("%d-%s", CodeTypeRegister, telephone))

	return nil
}

func (s *UserService) Login(telephone, email, enPassword, verificationCode, deviceId string, osVersion protos.GrpcReq_OS) (user *models.User, token string, tErr *TError) {
	var isNeedVerify bool = false
	var err error
	isNeedVerify, user, err = s.accountLogin(telephone, email, enPassword, verificationCode, deviceId, osVersion)
	if err != nil {
		return nil, "", NewTError(http.StatusBadRequest, err)
	}
	if isNeedVerify {
		return nil, "", NewTError(http.StatusAuthorizedRequired, errmsg.ErrAccountSecurityVerification)
	}
	// gen and save token
	token = utils.NewToken([]byte(user.UserId))
	if err = myRedis.Set(fmt.Sprintf("TK-%s", user.UserId), token, ThreeDays); err != nil {
		logger.Errorf("redis save token error: %v", err)
		return nil, "", NewTError(http.StatusInternalServerError, err)
	}
	return user, token, nil
}

func (s *UserService) Logout(token string, isMandatoryLogout bool) (tErr *TError) {
	isValid, payload, err := utils.TokenVerify(token)
	if err != nil {
		logger.Errorf("logout by token error: %s", err.Error())
		return NewTError(http.StatusInternalServerError, err)
	}
	if !isValid {
		return NewTError(http.StatusBadRequest, errmsg.ErrRepeatOperation)
	}

	key := fmt.Sprintf("TK-%s", string(payload))
	_token := myRedis.Get(key)
	if _token == "" {
		return NewTError(http.StatusBadRequest, errmsg.ErrRepeatOperation)
	}
	// TODO: if true, mandatory and will remove all online user
	if isMandatoryLogout {
		myRedis.Del(key)
	} else {
		myRedis.Del(key)
	}

	return nil
}

func (s *UserService) UpdateUserInfo(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	resp, _ = utils.NewGRPCResp(http.StatusOK, nil, "")

	// will over write the token value at rpc interceptor
	userId := req.GetToken()

	var err error
	var updateUserInfoReq protos.UpdateUserInfoReq
	if err = utils.UnmarshalGRPCReq(req, &updateUserInfoReq); err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = err.Error()
		logger.Errorf(`unmarshal error: %v`, err)
		return
	}

	pbProfile := updateUserInfoReq.GetProfile()
	if pbProfile == nil {
		resp.Code = http.StatusBadRequest
		resp.Message = errmsg.ErrInvalidParameters.Error()
		return
	}

	userProfile := converters.ConvertProto2EntityForUserProfile(pbProfile)
	userProfile.UserId = userId

	originUserProfile, err := us.userRepo.FindByUserId(userId)
	if err == errmsg.ErrInvalidUserId {
		resp.Code = http.StatusBadRequest
		resp.Message = errmsg.ErrInvalidUserId.Error()
		return
	}
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = errmsg.ErrInvalidParameters.Error()
		return
	}
	// nothing change, don't need to update
	if isProfileNothing2Update(originUserProfile.UserProfile, userProfile) {
		resp.Message = "profile unchanged"
		return
	}

	criteria := map[string]interface{}{
		"UserId": userId,
	}

	updateMap := utils.TransformStructToMap(userProfile)
	utils.RemoveMapProperties(updateMap, "UserId", "Telephone", "Email", "Avatar")

	err = us.userRepo.FindOneAndUpdateProfile(criteria, updateMap)
	if err != nil {
		logger.Errorf("[get user info] response marshal message error: %s", err.Error())

		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		return
	}
	return
}

func (s *UserService) ResetPassword(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	resp, _ = utils.NewGRPCResp(http.StatusOK, nil, "")

	var err error
	var resetPwdReq protos.ResetPasswordReq
	if err = utils.UnmarshalGRPCReq(req, &resetPwdReq); err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = err.Error()
		logger.Errorf(`unmarshal error: %v`, err)
		return
	}
	verificationCode := strings.Trim(resetPwdReq.GetVerificationCode(), "")
	oldPassword := strings.Trim(resetPwdReq.GetOldPassword(), "")
	newPassword := strings.Trim(resetPwdReq.GetNewPassword(), "")
	telephone := strings.Trim(resetPwdReq.GetTelephone(), "")
	email := strings.Trim(resetPwdReq.GetEmail(), "")

	if err = resetPwdParameterCalibration(verificationCode, oldPassword, newPassword, telephone, email); err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = err.Error()
		return
	}

	isRegistered, err := us.userRepo.IsTelephoneOrEmailRegistered(telephone, email)
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = err.Error()
		logger.Errorf("reset password error: %s", err.Error())
		return
	}
	if !isRegistered {
		resp.Code = http.StatusBadRequest
		resp.Message = errmsg.ErrUserNotExists.Error()
		logger.Warn("reset password failed: account not exist")
		return
	}

	if verificationCode != "" {
		// 1. code + newPassword
		code := myRedis.Get(fmt.Sprintf("%d-%s", CodeTypeResetPassword, telephone))
		if verificationCode != string(code) {
			resp.Code = http.StatusBadRequest
			resp.Message = "invalid verification code"
			return
		}
	} else {
		// 2. oldPassword + newPassword
		if oldPassword == newPassword {
			resp.Code = http.StatusBadRequest
			resp.Message = "the new password cannot be the same as the old one"
			return
		}
		enPassword := utils.NewSHA256(oldPassword, config.GetApiSecretKey())
		if telephone != "" {
			_, err = us.userRepo.QueryByTelephoneAndPassword(telephone, enPassword)
		} else {
			_, err = us.userRepo.QueryByEmailAndPassword(email, enPassword)
		}
		if err != nil {
			resp.Code = http.StatusBadRequest
			resp.Message = err.Error()
			logger.Errorf("reset password error: %s", err.Error())
			return
		}
	}

	enNewPassword := utils.NewSHA256(newPassword, config.GetApiSecretKey())
	if telephone != "" {
		if err = us.userRepo.ResetPasswordByTelephone(telephone, enNewPassword); err != nil {
			resp.Code = http.StatusInternalServerError
			resp.Message = err.Error()
			logger.Errorf("reset password error: %s", err.Error())
			return
		}
		myRedis.Del(fmt.Sprintf("%d-%s", CodeTypeResetPassword, telephone))
	} else {
		if err = us.userRepo.ResetPasswordByEmail(email, enNewPassword); err != nil {
			resp.Code = http.StatusInternalServerError
			resp.Message = err.Error()
			logger.Errorf("reset password error: %s", err.Error())
			return
		}
	}
	resetPwdResp := &protos.ResetPasswordResp{
		IsNeedReLogin: true,
	}
	resp.Data, err = utils.MarshalMessageToAny(resetPwdResp)
	if err != nil {
		logger.Errorf("[reset password] response marshal message error: %s", err.Error())
	}
	resp.Message = "password reset successful"

	return
}

func (s *UserService) GetUserInfo(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	resp, _ = utils.NewGRPCResp(http.StatusOK, nil, "")

	var err error
	var getUserReq protos.GetUserInfoReq
	if err = utils.UnmarshalGRPCReq(req, &getUserReq); err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = err.Error()
		logger.Errorf(`unmarshal error: %v`, err)
		return
	}
	if utils.IsEmptyStrings(getUserReq.GetUserId()) {
		resp.Code = http.StatusBadRequest
		resp.Message = errmsg.ErrInvalidParameters.Error()
		return
	}

	user, err := us.userRepo.FindByUserId(getUserReq.GetUserId())

	if err != nil {
		if err == errmsg.ErrInvalidUserId {
			resp.Code = http.StatusBadRequest
		} else {
			resp.Code = http.StatusInternalServerError
		}
		resp.Message = err.Error()
		return
	}
	userInfoResp := &protos.GetUserInfoResp{
		Profile: converters.ConvertEntity2ProtoForUserProfile(&user.UserProfile),
	}

	resp.Data, err = utils.MarshalMessageToAny(userInfoResp)
	if err != nil {
		logger.Errorf("[get user info] response marshal message error: %s", err.Error())
	}
	return
}

func (s *UserService) QueryUserInfo(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	resp, _ = utils.NewGRPCResp(http.StatusOK, nil, "")

	var err error
	var queryUserInfoReq protos.QueryUserInfoReq
	if err = utils.UnmarshalGRPCReq(req, &queryUserInfoReq); err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = err.Error()
		logger.Errorf(`unmarshal error: %v`, err)
		return
	}

	telephone := queryUserInfoReq.GetTelephone()
	email := queryUserInfoReq.GetEmail()
	if utils.IsEmptyStrings(telephone, email) {
		resp.Code = http.StatusBadRequest
		resp.Message = errmsg.ErrInvalidParameters.Error()
		return
	}

	userCriteria := make(map[string]interface{})
	if utils.IsEmptyStrings(telephone) {
		userCriteria["Email"] = email
	} else {
		userCriteria["Telephone"] = telephone
	}
	user, err := us.userRepo.FindOneUser(userCriteria)

	if err != nil {
		if err == errmsg.ErrUserNotExists {
			resp.Code = http.StatusBadRequest
		} else {
			resp.Code = http.StatusInternalServerError
		}
		resp.Message = err.Error()
		return
	}
	userInfoResp := &protos.QueryUserInfoResp{
		Profile: converters.ConvertEntity2ProtoForUserProfile(&user.UserProfile),
	}

	resp.Data, err = utils.MarshalMessageToAny(userInfoResp)
	if err != nil {
		logger.Errorf("[get user info] response marshal message error: %s", err.Error())
	}
	return
}

func isProfileNothing2Update(originProfile models.UserProfile, newProfile models.UserProfile) bool {
	return utils.DeepEqual(originProfile, newProfile)
}

func isVerificationCodeValid(verificationCode, telephone string) (isValid bool, err error) {
	fmt.Println("--------", myRedis)
	fmt.Println("========", fmt.Sprintf("%d-%s", CodeTypeRegister, telephone))
	code := myRedis.Get(fmt.Sprintf("%d-%s", CodeTypeRegister, telephone))
	if verificationCode != string(code) {
		return false, nil
	}
	return true, nil
}

func loginByTelephone(s *UserService, telephone, enPassword, verificationCode string) (user *models.User, err error) {
	// 1. login by password
	if enPassword != "" {
		user, err = us.userRepo.QueryByTelephoneAndPassword(telephone, enPassword)
		if err != nil {
			if err != errmsg.ErrAccountOrPwdInvalid {
				logger.Errorf("login by telephone and passwor error: %s", err)
				return nil, err
			}
		} else {
			return user, err
		}
	}
	// 2. login by verification code
	if verificationCode != "" {
		codeKey := fmt.Sprintf("%d-%s", CodeTypeLogin, telephone)
		code := myRedis.Get(codeKey)
		if verificationCode == string(code) {
			criteria := map[string]interface{}{
				"Telephone": telephone,
			}
			if user, err = us.userRepo.FindOneUser(criteria); err != nil {
				logger.Errorf("find user by telephone error: %s", err.Error())
				return nil, err
			}
			myRedis.Del(codeKey)
			return user, nil
		}
	}
	return user, err
}

func loginByEmail(s *UserService, email, enPassword, verificationCode string) (user *models.User, err error) {
	// 1. login by password
	if enPassword != "" {
		user, err = us.userRepo.QueryByEmailAndPassword(email, enPassword)
		if err != nil {
			if err != errmsg.ErrAccountOrPwdInvalid {
				logger.Errorf("login by email and passwor error: %s", err)
				return nil, err
			}
		} else {
			return user, nil
		}

	}
	// 2. login by verification code
	if verificationCode != "" {
		codeKey := fmt.Sprintf("%d-%s", CodeTypeLogin, email)
		code := myRedis.Get(codeKey)
		if verificationCode == string(code) {
			criteria := map[string]interface{}{
				"Email": email,
			}
			if user, err = us.userRepo.FindOneUser(criteria); err != nil {
				logger.Errorf("find user by email error: %s", err.Error())
				return nil, err
			}
			myRedis.Del(codeKey)
			return user, nil
		}
	}
	return user, nil

}

func (s *UserService) accountLogin(telephone, email, enPassword, verificationCode, deviceId string, osVersion protos.GrpcReq_OS) (isNeedVerify bool, user *models.User, err error) {
	var isLoginByTelephone bool = false
	if telephone != "" {
		isLoginByTelephone = true
	}
	// 1. login by verification code
	if verificationCode != "" {
		user, err = s.loginWithVerificationCode(isLoginByTelephone, telephone, email, verificationCode)
		if err != nil {
			return false, nil, err
		}
		if isNeedToSMSVerify(deviceId, osVersion, user) && user != nil {
			var condition = map[string]interface{}{
				"UserId": user.UserId,
			}
			var profile = map[string]interface{}{
				"DeviceId":  deviceId,
				"OsVersion": osVersion,
			}
			if err = us.userRepo.FindOneAndUpdateProfile(condition, profile); err != nil {
				return false, nil, err
			}
			return false, user, err
		}
	}
	// 2. login by password
	if enPassword != "" {
		user, err = us.loginWithPassword(isLoginByTelephone, telephone, email, enPassword)
		if err != nil {
			return false, nil, err
		}
		return isNeedToSMSVerify(deviceId, osVersion, user), user, err
	}
	return
}

func (s *UserService) loginWithPassword(isTelephone bool, telephone, email, enPassword string) (user *models.User, err error) {
	if isTelephone {
		return us.userRepo.QueryByTelephoneAndPassword(telephone, enPassword)
	} else {
		return us.userRepo.QueryByEmailAndPassword(email, enPassword)
	}
}

func (s *UserService) loginWithVerificationCode(isTelephone bool, telephone, email, verificationCode string) (user *models.User, err error) {
	var codeKey string = ""
	var criteria = make(map[string]interface{})
	if isTelephone {
		codeKey = fmt.Sprintf("%d-%s", CodeTypeLogin, telephone)
		criteria = map[string]interface{}{
			"Telephone": telephone,
		}
	} else {
		codeKey = fmt.Sprintf("%d-%s", CodeTypeLogin, email)
		criteria = map[string]interface{}{
			"Email": email,
		}
	}
	code := myRedis.Get(codeKey)
	if verificationCode == string(code) {
		if user, err = us.userRepo.FindOneUser(criteria); err != nil {
			logger.Errorf("find user by account error: %s", err.Error())
			return nil, err
		}
		myRedis.Del(codeKey)
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
