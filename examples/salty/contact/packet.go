package contact

import (
	"context"
	protos "goim-pro/api/protos/salty"
	"goim-pro/pkg/logs"
	"goim-pro/pkg/utils"
)

var (
	UserId = "01E4QYBXD0BVQYMYBTEDTRF9A4" // 13631210008
	logger = logs.GetLogger("INFO")
)

func RequestContact(t protos.ContactServiceClient) {
	reqContactReq := &protos.RequestContactReq{
		UserId: "01E4QYGJT86K2PS8CDHVXS0G95", // 13631210001
		Reason: "你好，交个朋友！",
	}
	anyData, _ := utils.MarshalMessageToAny(reqContactReq)
	grpcReq := &protos.GrpcReq{
		DeviceId: "",
		Version:  "",
		Language: 0,
		Os:       0,
		Token:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJNREZGTkZGWlFsaEVNRUpXVVZsTldVSlVSVVJVVWtZNVFUUT0iLCJleHAiOjE1ODU5MDYxMjAsImlhdCI6MTU4NTY0NjkyMCwiaXNzIjoic2FsdHlfaW0ifQ.fEaGkPMpt4wALF2OhP83dLwu0vcbWReFn92yh-zF4GY",
		Data:     anyData,
	}

	tr, err := t.RequestContact(context.Background(), grpcReq)
	if err != nil {
		logger.Errorf("login error: %s", err.Error())
	} else {
		printResp(tr)
	}
}

func AcceptContact(t protos.ContactServiceClient) {
	acpContactReq := &protos.AcceptContactReq{
		UserId: UserId,
	}
	anyData, _ := utils.MarshalMessageToAny(acpContactReq)
	grpcReq := &protos.GrpcReq{
		DeviceId: "",
		Version:  "",
		Language: 0,
		Os:       0,
		Token:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJNREZGTkZGWlIwcFVPRFpMTWxCVE9FTkVTRlpZVXpCSE9UVT0iLCJleHAiOjE1ODU5MDYzMjYsImlhdCI6MTU4NTY0NzEyNiwiaXNzIjoic2FsdHlfaW0ifQ.Ld7RNW5PhyGXtxYqe9eGs79Da9mNQYa79hy6R6K638M",
		Data:     anyData,
	}

	tr, err := t.AcceptContact(context.Background(), grpcReq)
	if err != nil {
		logger.Errorf("accept contact error: %s", err.Error())
	} else {
		printResp(tr)
	}
}

func RefusedContact(t protos.ContactServiceClient) {
	refusedContactReq := &protos.RefusedContactReq{
		UserId:               UserId,
		Reason:               "我喜欢小姐姐",
	}
	anyData, _ := utils.MarshalMessageToAny(refusedContactReq)
	grpcReq := &protos.GrpcReq{
		DeviceId: "",
		Version:  "",
		Language: 0,
		Os:       0,
		Token:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJNREZGTkZGWlIwcFVPRFpMTWxCVE9FTkVTRlpZVXpCSE9UVT0iLCJleHAiOjE1ODU5MDYzMjYsImlhdCI6MTU4NTY0NzEyNiwiaXNzIjoic2FsdHlfaW0ifQ.Ld7RNW5PhyGXtxYqe9eGs79Da9mNQYa79hy6R6K638M",
		Data:     anyData,
	}

	tr, err := t.RefusedContact(context.Background(), grpcReq)
	if err != nil {
		logger.Errorf("refused contact request error: %s", err.Error())
	} else {
		printResp(tr)
	}
}

func UpdateRemarkInfo(t protos.ContactServiceClient) {
	updateRemarkInfoReq := &protos.UpdateRemarkInfoReq{
		UserId: UserId,
		RemarkInfo: &protos.ContactRemark{
			RemarkName:  "喜洋洋",
			Description: "He's a crazy boy.",
			Telephones:  []string{"1361231222", "1369990440"},
			Tags:        []string{"Friend", "Boy"},
		},
	}
	anyData, _ := utils.MarshalMessageToAny(updateRemarkInfoReq)
	grpcReq := &protos.GrpcReq{
		DeviceId: "",
		Version:  "",
		Language: 0,
		Os:       0,
		Token:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJNREZGTkZGWlIwcFVPRFpMTWxCVE9FTkVTRlpZVXpCSE9UVT0iLCJleHAiOjE1ODU5MDYzMjYsImlhdCI6MTU4NTY0NzEyNiwiaXNzIjoic2FsdHlfaW0ifQ.Ld7RNW5PhyGXtxYqe9eGs79Da9mNQYa79hy6R6K638M",
		Data:     anyData,
	}

	tr, err := t.UpdateRemarkInfo(context.Background(), grpcReq)
	if err != nil {
		logger.Errorf("update remark info error: %s", err.Error())
	} else {
		printResp(tr)
	}
}

func printResp(resp *protos.GrpcResp) {
	logger.Infof("[code]: %d", resp.GetCode())
	logger.Infof("[data]: %s", resp.GetData().GetValue())
	logger.Infof("[message]: %s", resp.GetMessage())
}
