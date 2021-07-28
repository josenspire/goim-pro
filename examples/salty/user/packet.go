package user

import (
	"context"
	"fmt"
	protos "goim-pro/api/protos/salty"
	"goim-pro/pkg/logs"
	"goim-pro/pkg/utils"
	"log"
)

var (
	UserId = "01EMK01VY8820C95MH3ZGN3JXQ"
	logger = logs.GetLogger("INFO")
)

func ObtainSMSCode(t protos.SMSServiceClient, codeType protos.SMSOperationType, index int) {
	tel := fmt.Sprintf("1363121000%d", index)
	fmt.Println(tel)
	smsReq := protos.ObtainTelephoneSMSCodeReq{
		OperationType: codeType,
		Telephone:     tel,
	}
	anyData, _ := utils.MarshalMessageToAny(&smsReq)
	grpcReq := &protos.GrpcReq{
		DeviceId: "One Plus 7 Pro",
		Version:  "1",
		Language: 0,
		Os:       0,
		Token:    "",
		Data:     anyData,
	}
	// 调用 gRPC 接口
	tr, err := t.ObtainTelephoneSMSCode(context.Background(), grpcReq)
	//tr, err := t.Register(context.Background(), grpcReq)
	if err != nil {
		log.Fatalf("could not greet: %v", err.Error())
	}
	printResp(tr)

	//smsReq := protos.ObtainSMSCodeReq{
	//	CodeType:  protos.ObtainSMSCodeReq_REGISTER,
	//	Telephone: "13631210000",
	//}
	//
	//// 调用 gRPC 接口
	//tr, err := t.ObtainSMSCode(context.Background(), &smsReq)
	////tr, err := t.Register(context.Background(), grpcReq)
	//if err != nil {
	//	log.Fatalf("could not greet: %v", err.Error())
	//}
	//
	//logger.Infof("[code]: %d", tr)
}

func VerifyCode(t protos.SMSServiceClient) {
	tel := "13631210003"
	fmt.Println(tel)
	smsReq := protos.VerifyTelephoneSMSCodeReq{
		OperationType: protos.SMSOperationType_LOGIN,
		Telephone:     tel,
	}
	anyData, _ := utils.MarshalMessageToAny(&smsReq)
	grpcReq := &protos.GrpcReq{
		DeviceId: "One Plus 7 Pro",
		Version:  "1",
		Language: 0,
		Os:       0,
		Token:    "",
		Data:     anyData,
	}
	// 调用 gRPC 接口
	tr, err := t.VerifyTelephoneSMSCode(context.Background(), grpcReq)
	if err != nil {
		log.Fatalf("could not greet: %v", err.Error())
	}
	printResp(tr)
}

func ResetPasswordByTelephone(t protos.UserServiceClient, channel string) {
	var resetPasswordReq *protos.ResetPasswordReq
	resetPasswordReq = &protos.ResetPasswordReq{
		NewPassword: "112233445566",
		TargetAccount: &protos.ResetPasswordReq_Telephone{
			Telephone: "13631210003",
		},
	}

	anyData, _ := utils.MarshalMessageToAny(resetPasswordReq)
	grpcReq := &protos.GrpcReq{
		DeviceId: "LOCAL_DEV_123",
		Version:  "",
		Language: 0,
		Os:       0,
		Token:    "",
		Data:     anyData,
	}

	tr, err := t.ResetPassword(context.Background(), grpcReq)
	if err != nil {
		logger.Errorf("reset password error: %s", err.Error())
	} else {
		printResp(tr)
	}
}

func Register(t protos.UserServiceClient) {
	registerReq := &protos.RegisterReq{
		Password:         "1234567890",
		Profile: &protos.UserProfile{
			Telephone:   "13631210002",
			Email:       "12345672@qq.com",
			Nickname:    "JAMES02",
			Avatar:      "https://www.baidu.com/avatar/header1.png",
			Description: "Never settle",
			Sex:         protos.UserProfile_MALE,
			Birthday:    utils.MakeTimestamp(),
			Location:    "CHINA-ZHA",
		},
	}
	anyData, _ := utils.MarshalMessageToAny(registerReq)
	grpcReq := &protos.GrpcReq{
		DeviceId: "LOCAL_DEV",
		Version:  "",
		Language: 0,
		Os:       0,
		Token:    "",
		Data:     anyData,
	}

	tr, err := t.Register(context.Background(), grpcReq)
	if err != nil {
		logger.Errorf("register error: %s", err.Error())
	} else {
		printResp(tr)
	}
}

func Login(t protos.UserServiceClient, typeStr string, index int) {
	var loginReq *protos.LoginReq

	tel := fmt.Sprintf("1363121000%d", index)
	if typeStr == "TELEPHONE" {
		loginReq = &protos.LoginReq{
			TargetAccount: &protos.LoginReq_Telephone{
				Telephone: tel,
			},
			Password: "1234567890",
		}
	} else {
		loginReq = &protos.LoginReq{
			TargetAccount: &protos.LoginReq_Email{
				Email: "12345@qq.com",
			},
			Password: "1234567890",
		}
	}

	anyData, _ := utils.MarshalMessageToAny(loginReq)
	grpcReq := &protos.GrpcReq{
		DeviceId: "LOCAL_DEV",
		Version:  "",
		Language: 0,
		Os:       0,
		Token:    "",
		Data:     anyData,
	}

	tr, err := t.Login(context.Background(), grpcReq)
	if err != nil {
		logger.Errorf("login error: %s", err.Error())
	} else {
		printResp(tr)
	}
}

func LoginWithCode(t protos.UserServiceClient, typeStr string) {
	var loginReq *protos.LoginReq

	if typeStr == "TELEPHONE" {
		loginReq = &protos.LoginReq{
			TargetAccount: &protos.LoginReq_Telephone{
				Telephone: "13631210003",
			},
			Password:         "",
		}
	} else {
		loginReq = &protos.LoginReq{
			TargetAccount: &protos.LoginReq_Email{
				Email: "12345@qq.com",
			},
			Password: "1234567890",
		}
	}

	anyData, _ := utils.MarshalMessageToAny(loginReq)
	grpcReq := &protos.GrpcReq{
		DeviceId: "LOCAL_DEV",
		Version:  "",
		Language: 0,
		Os:       0,
		Token:    "",
		Data:     anyData,
	}

	tr, err := t.Login(context.Background(), grpcReq)
	if err != nil {
		logger.Errorf("login error: %s", err.Error())
	} else {
		printResp(tr)
	}
}

func Logout(t protos.UserServiceClient) {
	var logoutReq *protos.LogoutReq

	logoutReq = &protos.LogoutReq{
		IsMandatoryLogout: true,
	}

	anyData, _ := utils.MarshalMessageToAny(logoutReq)
	grpcReq := &protos.GrpcReq{
		DeviceId: "LOCAL_DEV",
		Version:  "",
		Language: 0,
		Os:       0,
		Token:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJNREZGTWtwWVRVTTVPRk5hV0UxSFJVZFdWRVJGUTFORU56Zz0iLCJleHAiOjE1ODM2NzYyNTYsImlhdCI6MTU4MzQxNzA1NiwiaXNzIjoic2FsdHlfaW0ifQ.TW40l918nMLITbHD4shmGHUOomlw2WC-SyFdimnG-cE",
		Data:     anyData,
	}

	tr, err := t.Logout(context.Background(), grpcReq)
	if err != nil {
		logger.Errorf("login error: %s", err.Error())
	} else {
		printResp(tr)
	}
}

func GetUserInfo(t protos.UserServiceClient) {
	var getUserInfoReq *protos.GetUserInfoReq

	getUserInfoReq = &protos.GetUserInfoReq{
		UserId: UserId,
	}

	anyData, _ := utils.MarshalMessageToAny(getUserInfoReq)
	grpcReq := &protos.GrpcReq{
		DeviceId: "LOCAL_DEV",
		Version:  "",
		Language: 0,
		Os:       0,
		Token:    "1234567890",
		Data:     anyData,
	}

	tr, err := t.GetUserInfo(context.Background(), grpcReq)
	if err != nil {
		logger.Errorf("login error: %s", err.Error())
	} else {
		printResp(tr)
	}
}

func QueryUserInfo(t protos.UserServiceClient, typeStr string) {
	var queryUserInfoReq *protos.QueryUserInfoReq
	switch typeStr {
	case "TELEPHONE":
		queryUserInfoReq = &protos.QueryUserInfoReq{
			TargetAccount: &protos.QueryUserInfoReq_Telephone{
				Telephone: "13631210003",
			},
		}
		break
	case "EMAIL":
		queryUserInfoReq = &protos.QueryUserInfoReq{
			TargetAccount: &protos.QueryUserInfoReq_Email{
				Email: "12345673@qq.com",
			},
		}
		break
	}

	anyData, _ := utils.MarshalMessageToAny(queryUserInfoReq)
	grpcReq := &protos.GrpcReq{
		DeviceId: "LOCAL_DEV",
		Version:  "",
		Language: 0,
		Os:       0,
		Token:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJNREZGVFVzd01WWlpPRGd5TUVNNU5VMUlNMXBIVGpOS1dGRT0iLCJleHAiOjE2MDc3NjE0NDIsImlhdCI6MTYwNzUwMjI0MiwiaXNzIjoic2FsdHlfaW0ifQ.FUrlMaANtx5uJMcRMqFnYbmohpVQWcI5QEHtMsccKKg",
		Data:     anyData,
	}

	tr, err := t.QueryUserInfo(context.Background(), grpcReq)
	if err != nil {
		logger.Errorf("login error: %s", err.Error())
	} else {
		printResp(tr)
	}
}

func UpdateUserInfo(t protos.UserServiceClient) {
	var updateUserInfoReq *protos.UpdateUserInfoReq
	updateUserInfoReq = &protos.UpdateUserInfoReq{
		Profile: &protos.UserProfile{
			UserId:      UserId,
			Telephone:   "13631210003",
			Email:       "12345@qq.com",
			Nickname:    "JAMES00012",
			Avatar:      "http://test/123456.jpg",
			Description: "de ma xiya",
			Sex:         1,
			Birthday:    0,
			Location:    "ZHA-CH",
		},
	}
	anyData, _ := utils.MarshalMessageToAny(updateUserInfoReq)
	grpcReq := &protos.GrpcReq{
		DeviceId: "LOCAL_DEV",
		Version:  "",
		Language: 0,
		Os:       0,
		Token:    "1234567890",
		Data:     anyData,
	}

	tr, err := t.UpdateUserInfo(context.Background(), grpcReq)
	if err != nil {
		logger.Errorf("login error: %s", err.Error())
	} else {
		printResp(tr)
	}
}

func printResp(resp *protos.GrpcResp) {
	logger.Infof("[code]: %d", resp.GetCode())
	logger.Infof("[data]: %s", resp.GetData().GetValue())
	logger.Infof("[message]: %s", resp.GetMessage())
}
