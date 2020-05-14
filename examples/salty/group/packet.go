package group

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

func CreateGroup(t protos.GroupServiceClient) {
	createGroupReq := &protos.CreateGroupReq{
		GroupName: "NEW_GROUP",
		MemberUserIdArr: []string{
			"MEMBER_001",
			"MEMBER_002",
		},
	}
	anyData, _ := utils.MarshalMessageToAny(createGroupReq)
	grpcReq := &protos.GrpcReq{
		DeviceId: "",
		Version:  "",
		Language: 0,
		Os:       0,
		Token:    "1234567890",
		Data:     anyData,
	}

	tr, err := t.CreateGroup(context.Background(), grpcReq)
	if err != nil {
		logger.Errorf("create group error: %s", err.Error())
	} else {
		printResp(tr)
	}
}

func printResp(resp *protos.GrpcResp) {
	logger.Infof("[code]: %d", resp.GetCode())
	logger.Infof("[data]: %s", resp.GetData().GetValue())
	logger.Infof("[message]: %s", resp.GetMessage())
}
