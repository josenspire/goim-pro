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
	"net/http"
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
	resp, _ = utils.NewGRPCResp(http.StatusOK, nil, "")

	var registerReq protos.RegisterReq
	if err = utils.UnmarshalGRPCReq(req, &registerReq); err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = err.Error()
		logger.Errorf(`unmarshal error: %v`, err)
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
	isRegistered, err := us.userRepo.IsTelephoneRegistered(userProfile.GetTelephone())
	if err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Message = fmt.Sprintf("server error, %s", err.Error())
		logger.Errorf("checking telephone validity error: %v", err.Error())
		return
	}
	if isRegistered {
		resp.Code = http.StatusBadRequest
		resp.Message = "this telephone has been registered, please login"
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
