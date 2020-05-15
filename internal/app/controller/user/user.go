package user

import (
	"context"
	protos "goim-pro/api/protos/salty"
	"goim-pro/config"
	"goim-pro/internal/app/services/converters"
	usersrv "goim-pro/internal/app/services/user"
	errmsg "goim-pro/pkg/errors"
	"goim-pro/pkg/http"
	"goim-pro/pkg/logs"
	"goim-pro/pkg/utils"
	"strings"
)

var (
	logger      = logs.GetLogger("INFO")
	userService *usersrv.UserService
)

type userServer struct {
}

func New() protos.UserServiceServer {
	userService = usersrv.New()
	return &userServer{}
}

func (u *userServer) Register(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
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

	userProfile := registerReq.GetProfile()
	verificationCode := registerReq.GetVerificationCode()
	password := registerReq.Password

	tErr := userService.Register(userProfile, password, verificationCode)
	if tErr != nil {
		resp.Code = tErr.Code
		resp.Message = tErr.Detail
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

func (u *userServer) Login(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
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
	telephone := loginReq.GetTelephone()
	email := loginReq.GetEmail()
	password := loginReq.GetPassword()
	verificationCode := loginReq.GetVerificationCode()
	deviceId := req.DeviceId
	osVersion := req.Os

	var enPassword string = ""
	if password != "" {
		enPassword = utils.NewSHA256(password, config.GetApiSecretKey())
	}

	user, token, tErr := userService.Login(telephone, email, enPassword, verificationCode, deviceId, osVersion)
	if tErr != nil {
		resp.Code = tErr.Code
		resp.Message = tErr.Detail
		return
	}

	loginResp := &protos.LoginResp{
		Token:   token,
		Profile: converters.ConvertEntity2ProtoForUserProfile(&user.UserProfile),
	}
	resp.Data, err = utils.MarshalMessageToAny(loginResp)
	if err != nil {
		logger.Errorf("login response marshal message error: %s", err.Error())
	}
	resp.Message = "user login successful"

	return
}

func (u *userServer) Logout(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	resp, _ = utils.NewGRPCResp(http.StatusOK, nil, "")

	var err error
	var logoutReq protos.LogoutReq
	if err = utils.UnmarshalGRPCReq(req, &logoutReq); err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = err.Error()
		logger.Errorf(`unmarshal error: %v`, err)
		return
	}

	token := req.GetToken()
	isMandatoryLogout := logoutReq.IsMandatoryLogout
	if tErr := userService.Logout(token, isMandatoryLogout); tErr != nil {
		resp.Code = tErr.Code
		resp.Message = tErr.Detail
		return
	}
	return
}

func (u *userServer) UpdateUserInfo(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	panic("implement me")
}

func (u *userServer) ResetPassword(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	panic("implement me")
}

func (u *userServer) GetUserInfo(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	panic("implement me")
}

func (u *userServer) QueryUserInfo(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	panic("implement me")
}

func registerParameterCalibration(req protos.RegisterReq) (err error) {
	csErr := errmsg.ErrInvalidParameters
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
	csErr := errmsg.ErrInvalidParameters
	req.VerificationCode = strings.Trim(req.GetVerificationCode(), "")
	req.Password = strings.Trim(req.GetPassword(), "")

	if utils.IsEmptyStrings(req.GetTelephone(), req.GetEmail()) || utils.IsEmptyStrings(req.GetPassword(), req.GetVerificationCode()) {
		err = csErr
	}
	return
}

func resetPwdParameterCalibration(verificationCode, oldPassword, newPassword string, telephone, email string) (err error) {
	csErr := errmsg.ErrInvalidParameters

	if utils.IsEmptyStrings(verificationCode, newPassword) || utils.IsEmptyStrings(oldPassword, newPassword) {
		err = csErr
	} else {
		if utils.IsEmptyStrings(telephone, email) {
			err = csErr
		}
	}
	return
}
