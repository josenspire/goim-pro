package usersrv

import (
	"context"
	"fmt"
	protos "goim-pro/api/protos/salty"
	"goim-pro/internal/app/repos"
	"goim-pro/internal/app/repos/user"
	"goim-pro/internal/app/services/converters"
	"goim-pro/pkg/logs"
	"goim-pro/pkg/utils"
)

type userService struct {
	userRepo user.IUserRepo
}

var logger = logs.GetLogger("INFO")

func New() protos.UserServiceServer {
	repoServer := repos.New()
	return &userService{
		userRepo: repoServer.UserRepo,
	}
}

func (us *userService) Register(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, err error) {
	resp = utils.NewResp(200, nil, "")

	var registerReq protos.RegisterReq
	err = utils.NewReq(req, &registerReq)
	if err != nil {
		resp.Code = 500
		logger.Errorf(`unmarshal error: %v`, err)
		return
	}
	isValid, err := isVerificationCodeValid(registerReq.GetRegisterType(), registerReq.GetVerificationCode())
	if !isValid {
		resp.Code = 400
		resp.Message = "verification code is invalid"
		logger.Warnf("verification code is invalid: %s", registerReq.GetVerificationCode())
		return
	}
	userProfile := registerReq.GetUserProfile()
	isRegistered, err := us.userRepo.IsTelephoneRegistered(userProfile.GetTelephone())
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
	err = us.userRepo.Register(&user.User{
		Password:    registerReq.GetPassword(),
		UserProfile: converters.ConvertRegisterUserProfile(userProfile),
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

func isVerificationCodeValid(registerType protos.RegisterReq_RegisterType, verificationCode string) (isValid bool, err error) {
	// TODO: should query from db
	if verificationCode != "123456" {
		return false, nil
	}
	return true, nil
}

func (us *userService) Login(ctx context.Context, req *protos.GrpcReq) (*protos.GrpcResp, error) {
	panic("implement me")
}
