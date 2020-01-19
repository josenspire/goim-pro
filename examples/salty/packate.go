package main

import (
	"context"
	protos "goim-pro/api/protos/salty"
	"goim-pro/pkg/utils"
	"log"
)

func obtainSMSCode(t protos.SMSServiceClient) {
	smsReq := protos.SMSReq{
		CodeType: protos.SMSReq_REGISTER,
		TargetAccount: &protos.SMSReq_Telephone{
			Telephone: "13631210000",
		},
	}
	anyData, _ := utils.MarshalMessageToAny(&smsReq)
	gprcReq := &protos.GrpcReq{
		Data: anyData,
	}
	// 调用 gRPC 接口
	tr, err := t.ObtainSMSCode(context.Background(), gprcReq)
	//tr, err := t.Register(context.Background(), gprcReq)
	if err != nil {
		log.Fatalf("could not greet: %v", err.Error())
	}
	printResp(tr)
}

func register(t protos.UserServiceClient) {
	registerReq := &protos.RegisterReq{
		RegisterType:     protos.RegisterReq_TELEPHONE,
		Password:         "1234567890",
		VerificationCode: "123456",
		UserProfile: &protos.UserProfile{
			Telephone:   "13631210003",
			Email:       "12345@qq.com",
			Username:    "JAMES001",
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
		DeviceID: "",
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
			LoginType: protos.LoginReq_TELEPHONE,
			TargetAccount: &protos.LoginReq_Telephone{
				Telephone: "13631210003",
			},
			Password: "1234567890",
		}
	} else {
		loginReq = &protos.LoginReq{
			LoginType: protos.LoginReq_EMAIL,
			TargetAccount: &protos.LoginReq_Email{
				Email: "12345@qq.com",
			},
			Password: "1234567890",
		}
	}

	anyData, _ := utils.MarshalMessageToAny(loginReq)
	grpcReq := &protos.GrpcReq{
		DeviceID: "",
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

func printResp(resp *protos.GrpcResp) {
	logger.Infof("[code]: %d", resp.GetCode())
	logger.Infof("[data]: %s", resp.GetData().GetValue())
	logger.Infof("[message]: %s", resp.GetMessage())
}
