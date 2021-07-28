package user

import (
	"context"
	protos "goim-pro/api/protos/salty"
	"goim-pro/config"
	"goim-pro/internal/app/services/converters"
	usersrv "goim-pro/internal/app/services/user"
	errmsg "goim-pro/pkg/errors"
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
	resp, _ = utils.NewGRPCResp(protos.StatusCode_STATUS_OK, nil, "")

	var err error
	var registerReq protos.RegisterReq
	if err = utils.UnmarshalGRPCReq(req, &registerReq); err != nil {
		resp.Code = protos.StatusCode_STATUS_BAD_REQUEST
		resp.Message = err.Error()
		logger.Errorf(`unmarshal error: %v`, err)
		return
	}
	if err = registerParameterCalibration(registerReq); err != nil {
		resp.Code = protos.StatusCode_STATUS_BAD_REQUEST
		resp.Message = err.Error()
		return
	}

	var (
		userProfile = registerReq.GetProfile()
		password = registerReq.Password
		deviceId = req.DeviceId
		osVersion = req.Os
	)

	tErr := userService.Register(userProfile, password, deviceId, osVersion)
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
	resp, _ = utils.NewGRPCResp(protos.StatusCode_STATUS_OK, nil, "")

	var err error
	var loginReq protos.LoginReq
	if err = utils.UnmarshalGRPCReq(req, &loginReq); err != nil {
		resp.Code = protos.StatusCode_STATUS_BAD_REQUEST
		resp.Message = err.Error()
		logger.Errorf(`unmarshal error: %v`, err)
		return
	}

	if err = loginParameterCalibration(&loginReq); err != nil {
		resp.Code = protos.StatusCode_STATUS_BAD_REQUEST
		resp.Message = err.Error()
		return
	}

	telephone := strings.Trim(loginReq.GetTelephone(), "")
	email := strings.Trim(loginReq.GetEmail(), "")
	password := strings.Trim(loginReq.GetPassword(), "")
	deviceId := strings.Trim(req.GetDeviceId(), "")
	osVersion := req.GetOs()

	var enPassword string = ""
	if password != "" {
		enPassword = utils.NewSHA256(password, config.GetApiSecretKey())
	}

	user, token, tErr := userService.Login(telephone, email, enPassword, deviceId, osVersion)
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
	resp, _ = utils.NewGRPCResp(protos.StatusCode_STATUS_OK, nil, "")

	var err error
	var logoutReq protos.LogoutReq
	if err = utils.UnmarshalGRPCReq(req, &logoutReq); err != nil {
		resp.Code = protos.StatusCode_STATUS_BAD_REQUEST
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
	resp, _ = utils.NewGRPCResp(protos.StatusCode_STATUS_OK, nil, "")

	var err error
	var updateUserInfoReq protos.UpdateUserInfoReq
	if err = utils.UnmarshalGRPCReq(req, &updateUserInfoReq); err != nil {
		resp.Code = protos.StatusCode_STATUS_BAD_REQUEST
		resp.Message = err.Error()
		logger.Errorf(`unmarshal error: %v`, err)
		return
	}

	// will over write the token value at rpc interceptor
	userId := req.GetToken()
	pbProfile := updateUserInfoReq.GetProfile()
	if pbProfile == nil {
		resp.Code = protos.StatusCode_STATUS_BAD_REQUEST
		resp.Message = errmsg.ErrInvalidParameters.Error()
		return
	}
	userProfile := converters.ConvertProto2EntityForUserProfile(pbProfile)
	userProfile.UserId = userId

	if tErr := userService.UpdateUserInfo(userId, &userProfile); tErr != nil {
		resp.Code = tErr.Code
		resp.Message = tErr.Detail
		return
	}
	return
}

func (u *userServer) ResetPassword(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	resp, _ = utils.NewGRPCResp(protos.StatusCode_STATUS_OK, nil, "")

	var err error
	var resetPwdReq protos.ResetPasswordReq
	if err = utils.UnmarshalGRPCReq(req, &resetPwdReq); err != nil {
		resp.Code = protos.StatusCode_STATUS_BAD_REQUEST
		resp.Message = err.Error()
		logger.Errorf(`unmarshal error: %v`, err)
		return
	}
	newPassword := strings.Trim(resetPwdReq.GetNewPassword(), "")
	telephone := strings.Trim(resetPwdReq.GetTelephone(), "")
	email := strings.Trim(resetPwdReq.GetEmail(), "")

	if err = resetPwdParameterCalibration(newPassword, telephone, email); err != nil {
		resp.Code = protos.StatusCode_STATUS_BAD_REQUEST
		resp.Message = err.Error()
		return
	}

	var (
		deviceId = req.DeviceId
		osVersion = req.Os
	)
	tErr := userService.ResetPassword(telephone, email, newPassword, deviceId, osVersion)
	if tErr != nil {
		resp.Code = tErr.Code
		resp.Message = tErr.Detail
		return
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

func (u *userServer) GetUserInfo(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	resp, _ = utils.NewGRPCResp(protos.StatusCode_STATUS_OK, nil, "")

	var err error
	var getUserReq protos.GetUserInfoReq
	if err = utils.UnmarshalGRPCReq(req, &getUserReq); err != nil {
		resp.Code = protos.StatusCode_STATUS_BAD_REQUEST
		resp.Message = err.Error()
		logger.Errorf(`unmarshal error: %v`, err)
		return
	}
	userId := strings.Trim(getUserReq.GetUserId(), "")
	if utils.IsEmptyStrings(userId) {
		resp.Code = protos.StatusCode_STATUS_BAD_REQUEST
		resp.Message = errmsg.ErrInvalidParameters.Error()
		return
	}

	userProfile, tErr := userService.GetUserInfo(userId)
	if tErr != nil {
		resp.Code = tErr.Code
		resp.Message = tErr.Detail
		return
	}

	userInfoResp := &protos.GetUserInfoResp{
		Profile: converters.ConvertEntity2ProtoForUserProfile(userProfile),
	}
	resp.Data, err = utils.MarshalMessageToAny(userInfoResp)
	if err != nil {
		logger.Errorf("[get user info] response marshal message error: %s", err.Error())
	}
	return
}

func (u *userServer) QueryUserInfo(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	resp, _ = utils.NewGRPCResp(protos.StatusCode_STATUS_OK, nil, "")

	var err error
	var queryUserInfoReq protos.QueryUserInfoReq
	if err = utils.UnmarshalGRPCReq(req, &queryUserInfoReq); err != nil {
		resp.Code = protos.StatusCode_STATUS_BAD_REQUEST
		resp.Message = err.Error()
		logger.Errorf(`unmarshal error: %v`, err)
		return
	}

	telephone := queryUserInfoReq.GetTelephone()
	email := queryUserInfoReq.GetEmail()
	if utils.IsEmptyStrings(telephone, email) {
		resp.Code = protos.StatusCode_STATUS_BAD_REQUEST
		resp.Message = errmsg.ErrInvalidParameters.Error()
		return
	}

	userProfile, tErr := userService.QueryUserInfo(telephone, email)
	if tErr != nil {
		resp.Code = tErr.Code
		resp.Message = tErr.Detail
		return
	}

	userInfoResp := &protos.QueryUserInfoResp{
		Profile: converters.ConvertEntity2ProtoForUserProfile(userProfile),
	}
	resp.Data, err = utils.MarshalMessageToAny(userInfoResp)
	if err != nil {
		logger.Errorf("[queryUserInfo] response marshal message error: %s", err.Error())
	}
	return
}

func registerParameterCalibration(req protos.RegisterReq) (err error) {
	csErr := errmsg.ErrInvalidParameters
	if utils.IsContainEmptyString(req.GetPassword()) || req.GetProfile() == nil {
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
	req.Password = strings.Trim(req.GetPassword(), "")

	if utils.IsEmptyStrings(req.GetTelephone(), req.GetEmail()) || utils.IsEmptyStrings(req.GetPassword()) {
		err = csErr
	}
	return
}

func resetPwdParameterCalibration(newPassword string, telephone, email string) (err error) {
	csErr := errmsg.ErrInvalidParameters

	if utils.IsEmptyStrings(newPassword) || utils.IsEmptyStrings(telephone, email) {
		err = csErr
	}
	return
}
