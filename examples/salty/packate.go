package main

import (
	"context"
	protos "goim-pro/api/protos/salty"
	"goim-pro/pkg/utils"
	"log"
)

var (
	UserId = "01E07SG858N3CGV5M1APVQKZYR"
)

func obtainSMSCode(t protos.SMSServiceClient, codeType protos.ObtainSMSCodeReq_CodeType) {
	smsReq := protos.ObtainSMSCodeReq{
		CodeType:  codeType,
		Telephone: "13631210000",
	}
	anyData, _ := utils.MarshalMessageToAny(&smsReq)
	gprcReq := &protos.GrpcReq{
		DeviceId: "One Plus 7 Pro",
		Version:  "1",
		Language: 0,
		Os:       0,
		Token:    "",
		Data:     anyData,
	}
	// 调用 gRPC 接口
	tr, err := t.ObtainSMSCode(context.Background(), gprcReq)
	//tr, err := t.Register(context.Background(), gprcReq)
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
	////tr, err := t.Register(context.Background(), gprcReq)
	//if err != nil {
	//	log.Fatalf("could not greet: %v", err.Error())
	//}
	//
	//logger.Infof("[code]: %d", tr)
}

func register(t protos.UserServiceClient) {
	registerReq := &protos.RegisterReq{
		Password:         "1234567890",
		VerificationCode: "123456",
		Profile: &protos.UserProfile{
			Telephone:   "13631210003",
			Email:       "12345@qq.com",
			Nickname:    "JAMES001",
			Avatar:      "https://www.baidu.com/avatar/header1.png",
			Description: "Never settle",
			Sex:         protos.UserProfile_MALE,
			Birthday:    utils.MakeTimestamp(),
			Location:    "CHINA-ZHA",
		},
	}
	anyData, _ := utils.MarshalMessageToAny(registerReq)
	grpcReq := &protos.GrpcReq{
		DeviceId: "",
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

func login(t protos.UserServiceClient, typeStr string) {
	var loginReq *protos.LoginReq

	if typeStr == "TELEPHONE" {
		loginReq = &protos.LoginReq{
			TargetAccount: &protos.LoginReq_Telephone{
				Telephone: "13631210003",
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
		DeviceId: "",
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

func getUserInfo(t protos.UserServiceClient) {
	var getUserInfoReq *protos.GetUserInfoReq

	getUserInfoReq = &protos.GetUserInfoReq{
		UserId: UserId,
	}

	anyData, _ := utils.MarshalMessageToAny(getUserInfoReq)
	grpcReq := &protos.GrpcReq{
		DeviceId: "",
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

func queryUserInfo(t protos.UserServiceClient, typeStr string) {
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
				Email: "12345@qq.com",
			},
		}
		break
	}

	anyData, _ := utils.MarshalMessageToAny(queryUserInfoReq)
	grpcReq := &protos.GrpcReq{
		DeviceId: "",
		Version:  "",
		Language: 0,
		Os:       0,
		Token:    "1234567890",
		Data:     anyData,
	}

	tr, err := t.QueryUserInfo(context.Background(), grpcReq)
	if err != nil {
		logger.Errorf("login error: %s", err.Error())
	} else {
		printResp(tr)
	}
}

func updateUserInfo(t protos.UserServiceClient) {
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
		DeviceId: "",
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
