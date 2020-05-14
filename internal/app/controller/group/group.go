package group

import (
	"context"
	protos "goim-pro/api/protos/salty"
	"goim-pro/internal/app/services/converters"
	groupsrv "goim-pro/internal/app/services/group"
	errmsg "goim-pro/pkg/errors"
	"goim-pro/pkg/http"
	"goim-pro/pkg/logs"
	"goim-pro/pkg/utils"
	"strings"
)

var (
	logger       = logs.GetLogger("INFO")
	groupService *groupsrv.GroupService
)

type groupServer struct {
}

func New() protos.GroupServiceServer {
	groupService = groupsrv.New()
	return &groupServer{}
}

func (s *groupServer) CreateGroup(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	resp, _ = utils.NewGRPCResp(http.StatusOK, nil, "")

	var err error
	var createGroupReq protos.CreateGroupReq
	if err = utils.UnmarshalGRPCReq(req, &createGroupReq); err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = err.Error()
		logger.Errorf(`unmarshal error: %v`, err)
		return
	}

	userId := req.GetToken()
	groupName := strings.Trim(createGroupReq.GroupName, "")
	memberIds := createGroupReq.MemberUserIdArr

	if len(memberIds) <= 0 || (len(memberIds) == 1 && memberIds[0] == userId) {
		resp.Code = http.StatusBadRequest
		resp.Message = errmsg.ErrIllegalOperation.Error()
		return
	}

	groupProfile, tErr := groupService.CreateGroup(userId, groupName, memberIds)
	if tErr != nil {
		resp.Code = tErr.Code
		resp.Message = tErr.Detail
		return
	}

	groupResp := &protos.CreateGroupResp{
		Profile: converters.ConvertEntity2ProtoForGroupProfile(groupProfile),
	}
	resp.Data, err = utils.MarshalMessageToAny(groupResp)
	if err != nil {
		logger.Errorf("create group response marshal message error: %s", err.Error())
	}
	resp.Message = "group create succeed"
	return
}

func (s *groupServer) JoinGroup(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	panic("implement me")
}

func (s *groupServer) QuitGroup(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	panic("implement me")
}

func (s *groupServer) AddGroupMember(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	panic("implement me")
}

func (s *groupServer) KickGroupMember(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	panic("implement me")
}

func (s *groupServer) UpdateGroupName(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	panic("implement me")
}

func (s *groupServer) UpdateGroupNotice(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	panic("implement me")
}

func (s *groupServer) UpdateMemberNickname(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	panic("implement me")
}
