package usersrv

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v7"
	protos "goim-pro/api/protos/salty"
	"goim-pro/config"
	"goim-pro/internal/app/constants"
	"goim-pro/internal/app/repos"
	. "goim-pro/internal/app/repos/user"
	"goim-pro/internal/app/services/converters"
	redsrv "goim-pro/pkg/db/redis"
	"goim-pro/pkg/logs"
	"goim-pro/pkg/utils"
	"net/http"
	"strings"
	"time"
)

var (
	logger = logs.GetLogger("INFO")
	crypto = utils.NewCrypto()

	expiresTime = time.Hour * time.Duration(24*3) // 3 days
	myRedis     *redis.Client
)

type userService struct {
	userRepo IUserRepo
}

func New() protos.UserServiceServer {
	myRedis = redsrv.NewRedisConnection().GetRedisClient()

	repoServer := repos.New()
	return &userService{
		userRepo: repoServer.UserRepo,
	}
}

func (us *userService) Register(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	resp, _ = utils.NewGRPCResp(http.StatusOK, nil, "")

	var err error
	var registerReq protos.RegisterReq
	if err = utils.UnmarshalGRPCReq(req, &registerReq); err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = err.Error()
		logger.Errorf(`unmarshal error: %v`, err)
		return
	}
	if err = registerParameterCalibration(registerReq); err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = err.Error()
		return
	}

	verificationCode := registerReq.GetVerificationCode()
	telephone := registerReq.GetProfile().GetTelephone()
	isValid, err := isVerificationCodeValid(verificationCode, telephone)
	if !isValid {
		resp.Code = http.StatusBadRequest
		resp.Message = "verification code is invalid"
		logger.Warnf("verification code is invalid: %s", verificationCode)
		return
	}
	userProfile := registerReq.GetProfile()
	userProfile.UserId = utils.NewULID()

	isRegistered, err := us.userRepo.IsTelephoneOrEmailRegistered(userProfile.GetTelephone(), userProfile.GetEmail())
	if err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Message = fmt.Sprintf("server error, %s", err.Error())
		logger.Errorf("checking telephone validity error: %v", err.Error())
		return
	}
	if isRegistered {
		resp.Code = http.StatusBadRequest
		resp.Message = "the telephone or email has been registered, please login"
		return
	}
	if err = us.userRepo.Register(&User{
		Password:    registerReq.GetPassword(),
		UserProfile: converters.ConvertProtoUserProfile2Entity(userProfile),
	}); err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		logger.Errorf("register user error: %v", err.Error())
		return
	}
	// remove cache
	myRedis.Del(fmt.Sprintf("%d-%s", constants.CodeTypeRegister, telephone))

	registerResp := &protos.RegisterResp{
		Profile: userProfile,
	}
	resp.Data, err = utils.MarshalMessageToAny(registerResp)
	if err != nil {
		logger.Errorf("register response marshal message error: %s", err.Error())
	}
	resp.Message = "user registration successful"
	return
}

func (us *userService) Login(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	resp, _ = utils.NewGRPCResp(http.StatusOK, nil, "")

	var err error
	var loginReq protos.LoginReq
	if err = utils.UnmarshalGRPCReq(req, &loginReq); err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = err.Error()
		logger.Errorf(`unmarshal error: %v`, err)
		return
	}
	if err = loginParameterCalibration(&loginReq); err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = err.Error()
		return
	}
	enPassword, err := crypto.AESEncrypt(loginReq.GetPassword(), config.GetApiSecretKey())
	if err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		return
	}
	var user *User
	if loginReq.GetTelephone() != "" {
		user, err = us.userRepo.QueryByTelephoneAndPassword(loginReq.GetTelephone(), enPassword)
	} else {
		user, err = us.userRepo.QueryByEmailAndPassword(loginReq.GetEmail(), enPassword)
	}
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = err.Error()
		logger.Errorf("user login error: %v", err.Error())
		return
	}

	// gen and save token
	token := utils.NewToken([]byte(user.UserId))
	if err = myRedis.Set(fmt.Sprintf("TK-%s", user.UserId), token, expiresTime).Err(); err != nil {
		logger.Errorf("redis save token error: %v", err)
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		return
	}

	loginResp := &protos.LoginResp{
		Token:   token,
		Profile: converters.ConvertProfileEntity2Proto(&user.UserProfile),
	}
	resp.Data, err = utils.MarshalMessageToAny(loginResp)
	if err != nil {
		logger.Errorf("login response marshal message error: %s", err.Error())
	}
	resp.Message = "user login successful"

	return
}

func (us *userService) Logout(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	resp, _ = utils.NewGRPCResp(http.StatusOK, nil, "")

	var err error
	var logoutReq protos.LogoutReq
	if err = utils.UnmarshalGRPCReq(req, &logoutReq); err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = err.Error()
		logger.Errorf(`unmarshal error: %v`, err)
		return
	}

	// if true, mandatory and will remove all online user
	if logoutReq.GetIsMandatoryLogout() {
		//TODO: remove all token from redis
	} else {
		//TODO: remove current token from redis
	}

	resp.Message = "user logout successful"

	return
}

func (us *userService) UpdateUserInfo(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
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
		resp.Message = utils.ErrInvalidParameters.Error()
		return
	}

	userProfile := converters.ConvertProtoUserProfile2Entity(pbProfile)
	userProfile.UserId = userId

	originUserProfile, err := us.userRepo.GetUserByUserId(userId)
	if err == utils.ErrInvalidUserId {
		resp.Code = http.StatusBadRequest
		resp.Message = utils.ErrInvalidUserId.Error()
		return
	}
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = utils.ErrInvalidParameters.Error()
		return
	}
	// nothing change, don't need to update
	if isProfileNothing2Update(originUserProfile.UserProfile, userProfile) {
		resp.Message = "profile unchanged"
		return
	}

	criteria := &User{}
	criteria.UserId = userId

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

func (us *userService) ResetPassword(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
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
		resp.Message = utils.ErrUserNotExists.Error()
		logger.Warn("reset password failed: account not exist")
		return
	}

	if verificationCode != "" {
		// 1. code + newPassword
		code, err := myRedis.Get(fmt.Sprintf("%d-%s", constants.CodeTypeResetPassword, telephone)).Bytes()
		if err != nil {
			logger.Errorf("read verification code error: %v", err)
			resp.Code = http.StatusInternalServerError
			resp.Message = err.Error()
			return
		}
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
		enPassword, err := crypto.AESEncrypt(oldPassword, config.GetApiSecretKey())
		if err != nil {
			resp.Code = http.StatusInternalServerError
			resp.Message = err.Error()
			return
		}
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

	enNewPassword, err := crypto.AESEncrypt(newPassword, config.GetApiSecretKey())
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = err.Error()
		logger.Errorf("reset password error: %s", err.Error())
		return
	}
	if telephone != "" {
		if err = us.userRepo.ResetPasswordByTelephone(telephone, enNewPassword); err != nil {
			resp.Code = http.StatusInternalServerError
			resp.Message = err.Error()
			logger.Errorf("reset password error: %s", err.Error())
			return
		}
		myRedis.Del(fmt.Sprintf("%d-%s", constants.CodeTypeResetPassword, telephone))
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

func (us *userService) GetUserInfo(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
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
		resp.Message = utils.ErrInvalidParameters.Error()
		return
	}

	user, err := us.userRepo.GetUserByUserId(getUserReq.GetUserId())

	if err != nil {
		if err == utils.ErrInvalidUserId {
			resp.Code = http.StatusBadRequest
		} else {
			resp.Code = http.StatusInternalServerError
		}
		resp.Message = err.Error()
		return
	}
	userInfoResp := &protos.GetUserInfoResp{
		Profile: converters.ConvertProfileEntity2Proto(&user.UserProfile),
	}

	resp.Data, err = utils.MarshalMessageToAny(userInfoResp)
	if err != nil {
		logger.Errorf("[get user info] response marshal message error: %s", err.Error())
	}
	return
}

func (us *userService) QueryUserInfo(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
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
		resp.Message = utils.ErrInvalidParameters.Error()
		return
	}

	userCriteria := &User{}
	if utils.IsEmptyStrings(telephone) {
		userCriteria.Email = email
	} else {
		userCriteria.Telephone = telephone
	}
	user, err := us.userRepo.FindOneUser(userCriteria)

	if err != nil {
		if err == utils.ErrUserNotExists {
			resp.Code = http.StatusBadRequest
		} else {
			resp.Code = http.StatusInternalServerError
		}
		resp.Message = err.Error()
		return
	}
	userInfoResp := &protos.QueryUserInfoResp{
		Profile: converters.ConvertProfileEntity2Proto(&user.UserProfile),
	}

	resp.Data, err = utils.MarshalMessageToAny(userInfoResp)
	if err != nil {
		logger.Errorf("[get user info] response marshal message error: %s", err.Error())
	}
	return
}

func isProfileNothing2Update(originProfile UserProfile, newProfile UserProfile) bool {
	return utils.DeepEqual(originProfile, newProfile)
}

func isVerificationCodeValid(verificationCode, telephone string) (isValid bool, err error) {
	code, err := myRedis.Get(fmt.Sprintf("%d-%s", constants.CodeTypeRegister, telephone)).Bytes()
	if err != nil {
		logger.Errorf("read verification code error: %v", err)
		return false, err
	}
	if verificationCode != string(code) {
		return false, nil
	}
	return true, nil
}

func registerParameterCalibration(req protos.RegisterReq) (err error) {
	csErr := errors.New("bad request, invalid parameters")
	if utils.IsContainEmptyString(req.GetPassword(), req.GetVerificationCode()) || req.GetProfile() == nil {
		err = csErr
		return
	}
	profile := req.GetProfile()
	//if registerType == protos.RegisterReq_TELEPHONE {
	//	if utils.IsContainEmptyString(profile.GetTelephone()) {
	//		err = csErr
	//	}
	//} else if registerType == protos.RegisterReq_EMAIL {
	//	if utils.IsContainEmptyString(profile.GetEmail()) {
	//		err = csErr
	//	}
	//}

	// telephone calibration
	if utils.IsContainEmptyString(profile.GetTelephone()) {
		err = csErr
	}
	return
}

func loginParameterCalibration(req *protos.LoginReq) (err error) {
	csErr := errors.New("bad request, invalid parameters")
	req.VerificationCode = strings.Trim(req.GetVerificationCode(), "")
	req.Password = strings.Trim(req.GetPassword(), "")

	if utils.IsContainEmptyString(req.GetPassword()) {
		err = csErr
	} else {
		if utils.IsEmptyStrings(req.GetTelephone(), req.GetEmail()) {
			err = csErr
		}
	}
	return
}

func resetPwdParameterCalibration(verificationCode, oldPassword, newPassword string, telephone, email string) (err error) {
	csErr := errors.New("bad request, invalid parameters")

	if utils.IsEmptyStrings(verificationCode, newPassword) || utils.IsEmptyStrings(oldPassword, newPassword) {
		err = csErr
	} else {
		if utils.IsEmptyStrings(telephone, email) {
			err = csErr
		}
	}
	return
}
