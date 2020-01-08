package usersrv

import (
	"context"
	"fmt"
	"goim-pro/api/protos"
	"goim-pro/internal/app/repos"
	"goim-pro/internal/app/repos/user"
	"goim-pro/internal/app/services/converts"
	"goim-pro/pkg/logs"
	"goim-pro/pkg/utils"
)

type userServer struct{}

var (
	logger      = logs.GetLogger("INFO")
	repoServer  = repos.New()
	userRepo    = repoServer.UserRepo
)

func New() protos.UserServiceServer {
	return &userServer{}
}

func (us *userServer) Register(ctx context.Context, req *protos.BasicClientRequest) (resp *protos.BasicServerResponse, err error) {
	resp = utils.NewResp(200, nil, "")

	var userReq protos.UserReq
	err = utils.NewReq(req, &userReq)
	if err != nil {
		resp.Code = 500
		logger.Errorf(`unmarshal error: %v`, err)
		return
	}
	isValid, err := isVerificationCodeValid(userReq.GetCodeType(), userReq.GetVerificationCode())
	if !isValid {
		resp.Code = 400
		resp.Message = "verification code is invalid"
		logger.Warnf("verification code is invalid: %s", userReq.GetVerificationCode())
		return
	}
	userProfile := userReq.GetUserProfile()
	isRegistered, err := userRepo.IsTelephoneRegistered(userProfile.GetTelephone())
	if err != nil {
		resp.Code = 500
		resp.Message = fmt.Sprintf("server error, %s", err.Error())
		logger.Errorf("checking telephone validity error: %v", err)
		return
	}
	if isRegistered {
		resp.Code = 400
		resp.Message = "this telephone has been registered, please login"
		return
	}
	err = userRepo.Register(&user.User{
		Password:    userReq.GetPassword(),
		UserProfile: converts.ConvertRegisterUserProfile(userProfile),
	})
	if err != nil {
		resp.Code = 500
		resp.Message = err.Error()
		logger.Errorf("register user error: %v", err)
	} else {
		resp.Message = "user registration successful"
	}
	return
}

func isVerificationCodeValid(codeType protos.CodeType, verificationCode string) (isValid bool, err error) {
	// TODO: should query from db
	if verificationCode != "123456" {
		return false, nil
	}
	return true, nil
}

func (us *userServer) Login(ctx context.Context, req *protos.BasicClientRequest) (*protos.BasicServerResponse, error) {
	panic("implement me")
}
