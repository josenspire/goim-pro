package groupsrv

import (
	"context"
	"github.com/jinzhu/gorm"
	protos "goim-pro/api/protos/salty"
	. "goim-pro/internal/app/constants"
	"goim-pro/internal/app/models"
	. "goim-pro/internal/app/repos/group"
	. "goim-pro/internal/app/repos/user"
	"goim-pro/internal/app/services/converters"
	mysqlsrv "goim-pro/pkg/db/mysql"
	"goim-pro/pkg/errors"
	"goim-pro/pkg/http"
	"goim-pro/pkg/logs"
	"goim-pro/pkg/utils"
	"strings"
)

var (
	logger  = logs.GetLogger("INFO")
	mysqlDB *gorm.DB
)

type groupService struct {
	userRepo  IUserRepo
	groupRepo IGroupRepo
}

func New() protos.GroupServiceServer {
	mysqlDB = mysqlsrv.NewMysqlConnection().GetMysqlInstance()

	//repoServer := repos.New(mysqlDB)
	return &groupService{
		userRepo:  NewUserRepo(mysqlDB),
		groupRepo: NewGroupRepo(mysqlDB),
	}
}

// CreateGroup - create new group
func (s *groupService) CreateGroup(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
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

	isOverflow, err := s.isGroupCountOverflow(userId)
	if err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		return
	}
	if isOverflow {
		resp.Code = http.StatusBadRequest
		resp.Message = errmsg.ErrGroupReachedLimit.Error()
		return
	}
	groupProfile := buildGroupProfile(userId, groupName, memberIds)
	ts := mysqlDB.Begin()
	groupProfile, err = s.groupRepo.CreateGroup(groupProfile)
	if err != nil {
		ts.Callback()
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		return
	}
	ts.Commit()

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

// JoinGroup - user join the group
func (s *groupService) JoinGroup(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
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
	// TODO:
	//requestReason := joinGroupReq.Reason

	if groupId == "" {
		resp.Code = http.StatusBadRequest
		resp.Message = errmsg.ErrIllegalOperation.Error()
		return
	}

	isMember, err := s.isGroupMember(groupId, userId)
	if err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		return
	}
	if isMember {
		resp.Code = http.StatusRepeatOperation
		resp.Message = errmsg.ErrRepeatedlyJoinGroup.Error()
		return
	}

	newMember := models.NewMember(userId, "")
	newMember.GroupId = groupId

	if err := s.groupRepo.InsertMembers(&newMember); err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		return
	}

	var condition = map[string]interface{}{
		"groupId": groupId,
	}
	groupProfile, err := s.groupRepo.FindOneGroup(condition)
	if err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
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

// QuitGroup - user quit the group
func (s *groupService) QuitGroup(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
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
	isMember, err := s.isGroupMember(groupId, userId)
	if err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		return
	}
	if !isMember {
		resp.Code = http.StatusRepeatOperation
		resp.Message = errmsg.ErrNotGroupMembers.Error()
		return
	}

	memberIds := []string{userId}
	count, err := s.groupRepo.RemoveMembers(groupId, memberIds, true)
	if err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		return
	}
	if count == 0 {
		resp.Code = http.StatusInternalServerError
		resp.Message = "server error, quit group fail"
		return
	}
	resp.Message = "quit group succeed"
	return
}

// AddGroupMember - add multiple members into group
func (s *groupService) AddGroupMember(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
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

	condition := map[string]interface{}{
		"groupId": groupId,
	}
	groupProfile, err := s.groupRepo.FindOneGroup(condition)
	if err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		return
	}
	if groupProfile == nil {
		resp.Code = http.StatusBadRequest
		resp.Message = errmsg.ErrGroupNotExists.Error()
		return
	}
	members := groupProfile.Members
	if isOutOfMemberLimit(len(members), len(memberIds)) {
		resp.Code = http.StatusBadRequest
		resp.Message = errmsg.ErrGroupMemberReachedLimit.Error()
		return
	}

	newMembers := buildMembersObject(groupId, memberIds)
	if err = s.groupRepo.InsertMembers(newMembers...); err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		return
	}
	newGroupProfile, err := s.groupRepo.FindOneGroup(condition)
	if err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		return
	}
	var addGroupMemberResp = &protos.AddGroupMemberResp{
		Profile: converters.ConvertEntity2ProtoForGroupProfile(newGroupProfile),
	}
	resp.Data, err = utils.MarshalMessageToAny(addGroupMemberResp)
	if err != nil {
		logger.Errorf("[addGroupMember] response marshal message error: %s", err.Error())
	}
	return
}

// AddGroupMember - kick one member from a group
func (s *groupService) KickGroupMember(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
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

	condition := map[string]interface{}{
		"groupId": groupId,
	}
	groupProfile, err := s.groupRepo.FindOneGroup(condition)
	if err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		return
	}
	if groupProfile == nil {
		resp.Code = http.StatusBadRequest
		resp.Message = errmsg.ErrGroupNotExists.Error()
		return
	}

	isMember, err := s.isGroupMember(groupId, userId)
	if err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		return
	}
	if !isMember {
		resp.Code = http.StatusRepeatOperation
		resp.Message = errmsg.ErrNotGroupMembers.Error()
		return
	}

	// check out currently user permission
	if isGroupManager(userId, groupProfile.OwnerUserId) {
		resp.Code = http.StatusRequestForbidden
		resp.Message = errmsg.ErrOperationForbidden.Error()
		return
	}

	_, err = s.groupRepo.RemoveMembers(groupId, []string{memberUserId}, true)
	if err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		return
	}

	newGroupProfile, err := s.groupRepo.FindOneGroup(condition)
	if err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		return
	}
	var kickGroupMemberResp = &protos.KickGroupMemberResp{
		Profile: converters.ConvertEntity2ProtoForGroupProfile(newGroupProfile),
	}
	resp.Data, err = utils.MarshalMessageToAny(kickGroupMemberResp)
	if err != nil {
		logger.Errorf("[kickGroupMemberResp] response marshal message error: %s", err.Error())
	}
	return
}

// UpdateGroupName - update group name by group manager
func (s *groupService) UpdateGroupName(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
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
	// handle empty group name to default
	if newGroupName == "" {
		newGroupName = models.DefaultGroupName
	}

	condition := map[string]interface{}{
		"groupId": groupId,
	}
	groupProfile, err := s.groupRepo.FindOneGroup(condition)
	if err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		return
	}
	if groupProfile == nil {
		resp.Code = http.StatusBadRequest
		resp.Message = errmsg.ErrGroupNotExists.Error()
		return
	}

	// check out user permission
	if !isGroupManager(userId, groupProfile.OwnerUserId) {
		resp.Code = http.StatusRequestForbidden
		resp.Message = errmsg.ErrOperationForbidden.Error()
		return
	}

	updated := map[string]interface{}{
		"name": newGroupName,
	}
	newGroupProfile, err := s.groupRepo.FindOneGroupAndUpdate(condition, updated)
	if err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		return
	}
	var updateGroupNameResp = &protos.UpdateGroupNameResp{
		Profile: converters.ConvertEntity2ProtoForGroupProfile(newGroupProfile),
	}
	resp.Data, err = utils.MarshalMessageToAny(updateGroupNameResp)
	if err != nil {
		logger.Errorf("[updateGroupNameResp] response marshal message error: %s", err.Error())
	}

	return
}

// UpdateGroupNotice - update group notice by group manager
func (s *groupService) UpdateGroupNotice(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
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

	condition := map[string]interface{}{
		"groupId": groupId,
	}
	groupProfile, err := s.groupRepo.FindOneGroup(condition)
	if err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		return
	}
	if groupProfile == nil {
		resp.Code = http.StatusBadRequest
		resp.Message = errmsg.ErrGroupNotExists.Error()
		return
	}

	// check out user permission
	if !isGroupManager(userId, groupProfile.OwnerUserId) {
		resp.Code = http.StatusRequestForbidden
		resp.Message = errmsg.ErrOperationForbidden.Error()
		return
	}

	updated := map[string]interface{}{
		"notice": newNotice,
	}
	newGroupProfile, err := s.groupRepo.FindOneGroupAndUpdate(condition, updated)
	if err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		return
	}
	var updateGroupNoticeResp = &protos.UpdateGroupNoticeResp{
		Profile: converters.ConvertEntity2ProtoForGroupProfile(newGroupProfile),
	}
	resp.Data, err = utils.MarshalMessageToAny(updateGroupNoticeResp)
	if err != nil {
		logger.Errorf("[updateGroupNoticeResp] response marshal message error: %s", err.Error())
	}
	return
}

// UpdateMemberNickname - update current user's alias in group
func (s *groupService) UpdateMemberNickname(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
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

	condition := map[string]interface{}{
		"groupId": groupId,
		"userId":  userId,
	}
	updated := map[string]interface{}{
		"alias": newAlias,
	}
	memberProfile, err := s.groupRepo.FindOneMemberAndUpdate(condition, updated)
	if err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		return
	}
	if memberProfile == nil {
		resp.Code = http.StatusBadRequest
		resp.Message = errmsg.ErrNotGroupMembers.Error()
		return
	}

	searchGroup := map[string]interface{}{
		"groupId": groupId,
	}
	groupProfile, err := s.groupRepo.FindOneGroup(searchGroup)
	if err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		return
	}
	if groupProfile == nil { // current stage profile would not be nil
		resp.Code = http.StatusInternalServerError
		resp.Message = errmsg.ErrSystemUncheckException.Error()
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

// isGroupMember - check target user whether is a member of the group
func (s *groupService) isGroupMember(groupId, userId string) (bool, error) {
	var condition = map[string]interface{}{
		"groupId": groupId,
		"userId":  userId,
	}
	memberProfile, err := s.groupRepo.FindOneMember(condition)
	if err != nil {
		return false, err
	}
	return memberProfile == nil, nil
}

// isGroupCountOverflow - check target user's group whether is reach limit count
func (s *groupService) isGroupCountOverflow(userId string) (isOverflow bool, err error) {
	condition := map[string]interface{}{
		"userId": userId,
	}
	totalNum, err := s.groupRepo.CountGroup(condition)
	if err != nil {
		return false, err
	}
	if totalNum >= MaximumNumberOfGroups {
		return true, nil
	}
	return false, nil
}

func buildGroupProfile(userId string, groupName string, memberIds []string) *models.Group {
	var members = make([]models.Member, len(memberIds)-1)
	for i, userId := range memberIds {
		members[i] = models.NewMember(userId, "")
	}
	return models.NewGroup(utils.NewULID(), userId, groupName, members)
}

func buildMembersObject(groupId string, memberIds []string) (members []*models.Member) {
	members = make([]*models.Member, len(memberIds)-1)
	for i, userId := range memberIds {
		members[i] = &models.Member{
			UserId:  userId,
			GroupId: groupId,
		}
	}
	return
}

func isOutOfMemberLimit(orgSize int, newSize int) bool {
	return MaximumNumberOfGroupMembers < orgSize+newSize
}

func isGroupManager(userId string, targetUserId string) bool {
	return strings.EqualFold(userId, targetUserId)
}
