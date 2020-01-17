package usersrv

import (
	"context"
	"errors"
	"fmt"
	protos "goim-pro/api/protos/salty"
	"goim-pro/config"
	"goim-pro/internal/app/repos"
	"goim-pro/internal/app/repos/user"
	"goim-pro/internal/app/services/converters"
	"goim-pro/pkg/logs"
	"goim-pro/pkg/utils"
	"net/http"
)

var logger = logs.GetLogger("INFO")
var crypto = utils.NewCrypto()

type userService struct {
	userRepo user.IUserRepo
}

func New() protos.UserServiceServer {
	repoServer := repos.New()
	return &userService{
		userRepo: repoServer.UserRepo,
	}
}

func (us *userService) Register(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, grpcErr error) {
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
	isValid, err := isVerificationCodeValid(registerReq.GetRegisterType(), registerReq.GetVerificationCode())
	if !isValid {
		resp.Code = http.StatusBadRequest
		resp.Message = "verification code is invalid"
		logger.Warnf("verification code is invalid: %s", registerReq.GetVerificationCode())
		return
	}
	userProfile := registerReq.GetUserProfile()
	userProfile.UserID = utils.NewULID()

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
	if err = us.userRepo.Register(&user.User{
		Password:    registerReq.GetPassword(),
		UserProfile: converters.ConvertRegisterUserProfile(userProfile),
	}); err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		logger.Errorf("register user error: %v", err.Error())
		return
	}
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

func (us *userService) Login(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, grpcErr error) {
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
	var user *user.User
	if loginReq.GetEmail() != "" {
		user, err = us.userRepo.LoginByTelephone(loginReq.GetTelephone(), enPassword)
	} else {
		user, err = us.userRepo.LoginByEmail(loginReq.GetEmail(), enPassword)
	}
	if err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		logger.Errorf("register user error: %v", err.Error())
		return
	}

	// TODO: should verify
	registerResp := &protos.RegisterResp{
		Profile: converters.ConvertLoginResp(user.UserProfile),
	}
	resp.Data, err = utils.MarshalMessageToAny(registerResp)
	if err != nil {
		logger.Errorf("register response marshal message error: %s", err.Error())
	}
	resp.Message = "user registration successful"

	return
}

func isVerificationCodeValid(registerType protos.RegisterReq_RegisterType, verificationCode string) (isValid bool, err error) {
	// TODO: should query from db
	if verificationCode != "123456" {
		return false, nil
	}
	return true, nil
}

func registerParameterCalibration(req protos.RegisterReq) (err error) {
	csErr := errors.New("bad request, invalid parameters")
	if utils.IsContainEmptyString(req.GetPassword(), req.GetVerificationCode()) || req.GetUserProfile() == nil {
		err = csErr
		return
	}
	profile := req.GetUserProfile()
	registerType := req.RegisterType
	if registerType == protos.RegisterReq_TELEPHONE {
		if utils.IsContainEmptyString(profile.GetTelephone()) {
			err = csErr
		}
	} else if registerType == protos.RegisterReq_EMAIL {
		if utils.IsContainEmptyString(profile.GetEmail()) {
			err = csErr
		}
	}
	return
}

func loginParameterCalibration(req *protos.LoginReq) (err error) {
	csErr := errors.New("bad request, invalid parameters")
	if utils.IsContainEmptyString(req.GetPassword(), req.GetEmail(), req.GetTelephone()) {
		err = csErr
	}
	return
}
