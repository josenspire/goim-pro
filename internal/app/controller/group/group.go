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
	resp, _ = utils.NewGRPCResp(http.StatusOK, nil, "")

	var err error
	var joinGroupReq protos.JoinGroupReq
	if err = utils.UnmarshalGRPCReq(req, &joinGroupReq); err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = err.Error()
		logger.Errorf(`unmarshal error: %v`, err)
		return
	}

	userId := req.GetToken()
	groupId := joinGroupReq.GroupId
	// TODO: should fill with the notification
	requestReason := joinGroupReq.Reason

	if groupId == "" {
		resp.Code = http.StatusBadRequest
		resp.Message = errmsg.ErrIllegalOperation.Error()
		return
	}

	groupProfile, tErr := groupService.JoinGroup(userId, groupId, requestReason)
	if tErr != nil {
		resp.Code = tErr.Code
		resp.Message = tErr.Detail
		return
	}

	joinGroupResp := &protos.JoinGroupResp{
		Profile: converters.ConvertEntity2ProtoForGroupProfile(groupProfile),
	}
	resp.Data, err = utils.MarshalMessageToAny(joinGroupResp)
	if err != nil {
		logger.Errorf("join group response marshal message error: %s", err.Error())
	}
	resp.Message = "join group succeed"
	return
}

func (s *groupServer) QuitGroup(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	resp, _ = utils.NewGRPCResp(http.StatusOK, nil, "")

	var err error
	var quitGroupReq protos.QuitGroupReq
	if err = utils.UnmarshalGRPCReq(req, &quitGroupReq); err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = err.Error()
		logger.Errorf(`unmarshal error: %v`, err)
		return
	}

	userId := req.GetToken()
	groupId := quitGroupReq.GroupId

	if groupId == "" {
		resp.Code = http.StatusBadRequest
		resp.Message = errmsg.ErrIllegalOperation.Error()
		return
	}

	if tErr := groupService.QuitGroup(userId, groupId); tErr != nil {
		resp.Code = tErr.Code
		resp.Message = tErr.Detail
		return
	}
	return
}

func (s *groupServer) AddGroupMember(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	resp, _ = utils.NewGRPCResp(http.StatusOK, nil, "")

	var err error
	var addGroupMemberReq protos.AddGroupMemberReq
	if err = utils.UnmarshalGRPCReq(req, &addGroupMemberReq); err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = err.Error()
		logger.Errorf(`unmarshal error: %v`, err)
		return
	}

	//userId := req.GetToken()
	groupId := addGroupMemberReq.GroupId
	memberIds := addGroupMemberReq.MemberUserIdArr

	if groupId == "" || len(memberIds) == 0 {
		resp.Code = http.StatusBadRequest
		resp.Message = errmsg.ErrIllegalOperation.Error()
		return
	}

	groupProfile, tErr := groupService.AddGroupMember(groupId, memberIds)
	if tErr != nil {
		resp.Code = tErr.Code
		resp.Message = tErr.Detail
		return
	}

	var addGroupMemberResp = &protos.AddGroupMemberResp{
		Profile: converters.ConvertEntity2ProtoForGroupProfile(groupProfile),
	}
	resp.Data, err = utils.MarshalMessageToAny(addGroupMemberResp)
	if err != nil {
		logger.Errorf("[addGroupMember] response marshal message error: %s", err.Error())
	}
	return
}

func (s *groupServer) KickGroupMember(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	resp, _ = utils.NewGRPCResp(http.StatusOK, nil, "")

	var err error
	var kickGroupMemberReq protos.KickGroupMemberReq
	if err = utils.UnmarshalGRPCReq(req, &kickGroupMemberReq); err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = err.Error()
		logger.Errorf(`unmarshal error: %v`, err)
		return
	}

	userId := req.GetToken()
	groupId := kickGroupMemberReq.GroupId
	memberUserId := kickGroupMemberReq.MemberUserId

	if utils.IsContainEmptyString(groupId, memberUserId) {
		resp.Code = http.StatusBadRequest
		resp.Message = errmsg.ErrIllegalOperation.Error()
		return
	}

	groupProfile, tErr := groupService.KickGroupMember(userId, groupId, memberUserId)
	if tErr != nil {
		resp.Code = tErr.Code
		resp.Message = tErr.Detail
		return
	}
	var kickGroupMemberResp = &protos.KickGroupMemberResp{
		Profile: converters.ConvertEntity2ProtoForGroupProfile(groupProfile),
	}
	resp.Data, err = utils.MarshalMessageToAny(kickGroupMemberResp)
	if err != nil {
		logger.Errorf("[kickGroupMemberResp] response marshal message error: %s", err.Error())
	}
	return
}

func (s *groupServer) UpdateGroupName(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	resp, _ = utils.NewGRPCResp(http.StatusOK, nil, "")

	var err error
	var updateGroupNameReq protos.UpdateGroupNameReq
	if err = utils.UnmarshalGRPCReq(req, &updateGroupNameReq); err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = err.Error()
		logger.Errorf(`unmarshal error: %v`, err)
		return
	}

	userId := req.GetToken()
	groupId := updateGroupNameReq.GroupId
	newGroupName := updateGroupNameReq.GroupName

	if groupId == "" {
		resp.Code = http.StatusBadRequest
		resp.Message = errmsg.ErrIllegalOperation.Error()
		return
	}

	groupProfile, tErr := groupService.UpdateGroupName(userId, groupId, newGroupName)
	if tErr != nil {
		resp.Code = tErr.Code
		resp.Message = tErr.Detail
		return
	}

	var updateGroupNameResp = &protos.UpdateGroupNameResp{
		Profile: converters.ConvertEntity2ProtoForGroupProfile(groupProfile),
	}
	resp.Data, err = utils.MarshalMessageToAny(updateGroupNameResp)
	if err != nil {
		logger.Errorf("[updateGroupNameResp] response marshal message error: %s", err.Error())
	}
	return
}

func (s *groupServer) UpdateGroupNotice(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	resp, _ = utils.NewGRPCResp(http.StatusOK, nil, "")

	var err error
	var updateGroupNoticeReq protos.UpdateGroupNoticeReq
	if err = utils.UnmarshalGRPCReq(req, &updateGroupNoticeReq); err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = err.Error()
		logger.Errorf(`unmarshal error: %v`, err)
		return
	}

	userId := req.GetToken()
	groupId := updateGroupNoticeReq.GroupId
	newNotice := updateGroupNoticeReq.Notice

	if groupId == "" {
		resp.Code = http.StatusBadRequest
		resp.Message = errmsg.ErrIllegalOperation.Error()
		return
	}
	groupProfile, tErr := groupService.UpdateGroupNotice(userId, groupId, newNotice)
	if tErr != nil {
		resp.Code = tErr.Code
		resp.Message = tErr.Detail
		return
	}
	var updateGroupNoticeResp = &protos.UpdateGroupNoticeResp{
		Profile: converters.ConvertEntity2ProtoForGroupProfile(groupProfile),
	}
	resp.Data, err = utils.MarshalMessageToAny(updateGroupNoticeResp)
	if err != nil {
		logger.Errorf("[updateGroupNoticeResp] response marshal message error: %s", err.Error())
	}
	return
}

func (s *groupServer) UpdateMemberNickname(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	resp, _ = utils.NewGRPCResp(http.StatusOK, nil, "")

	var err error
	var updateMemberNicknameReq protos.UpdateMemberNicknameReq
	if err = utils.UnmarshalGRPCReq(req, &updateMemberNicknameReq); err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = err.Error()
		logger.Errorf(`unmarshal error: %v`, err)
		return
	}

	userId := req.GetToken()
	groupId := updateMemberNicknameReq.GroupId
	newAlias := updateMemberNicknameReq.MemberNickname
	if groupId == "" {
		resp.Code = http.StatusBadRequest
		resp.Message = errmsg.ErrIllegalOperation.Error()
		return
	}

	groupProfile, tErr := groupService.UpdateMemberNickname(userId, groupId, newAlias)
	if tErr != nil {
		resp.Code = tErr.Code
		resp.Message = tErr.Detail
		return
	}
	var updateMemberNicknameResp = &protos.UpdateMemberNicknameResp{
		Profile: converters.ConvertEntity2ProtoForGroupProfile(groupProfile),
	}
	resp.Data, err = utils.MarshalMessageToAny(updateMemberNicknameResp)
	if err != nil {
		logger.Errorf("[updateGroupNoticeResp] response marshal message error: %s", err.Error())
	}
	return
}
