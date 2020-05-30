package groupsrv

import (
	"github.com/jinzhu/gorm"
	protos "goim-pro/api/protos/salty"
	. "goim-pro/internal/app/constants"
	"goim-pro/internal/app/models"
	. "goim-pro/internal/app/models/errors"
	. "goim-pro/internal/app/repos/group"
	. "goim-pro/internal/app/repos/user"
	mysqlsrv "goim-pro/pkg/db/mysql"
	"goim-pro/pkg/errors"
	"goim-pro/pkg/logs"
	"goim-pro/pkg/utils"
	"strings"
)

var (
	logger  = logs.GetLogger("INFO")
	mysqlDB *gorm.DB

	userRepo  IUserRepo
	groupRepo IGroupRepo
)

type GroupService struct{}

func New() *GroupService {
	mysqlDB = mysqlsrv.NewMysql()
	userRepo = NewUserRepo(mysqlDB)
	groupRepo = NewGroupRepo(mysqlDB)
	return &GroupService{}
}

// CreateGroup - create new group
func (s *GroupService) CreateGroup(userId, groupName string, memberIds []string) (profile *models.Group, tErr *TError) {
	if len(memberIds) <= 0 || (len(memberIds) == 1 && memberIds[0] == userId) {
		return nil, NewTError(protos.StatusCode_STATUS_BAD_REQUEST, errmsg.ErrIllegalOperation)
	}

	isOverflow, err := s.isGroupCountOverflow(userId)
	if err != nil {
		return nil, NewTError(protos.StatusCode_STATUS_INTERNAL_SERVER_ERROR, err)
	}
	if isOverflow {
		return nil, NewTError(protos.StatusCode_STATUS_BAD_REQUEST, errmsg.ErrGroupReachedLimit)
	}
	groupProfile := buildGroupProfile(userId, groupName, memberIds)
	ts := mysqlDB.Begin()
	groupProfile, err = groupRepo.CreateGroup(groupProfile)
	if err != nil {
		ts.Callback()
		return nil, NewTError(protos.StatusCode_STATUS_INTERNAL_SERVER_ERROR, err)
	}
	profile, err = groupRepo.FindOneGroup(map[string]interface{}{"groupId": groupProfile.GroupId})
	if err != nil {
		ts.Callback()
		return nil, NewTError(protos.StatusCode_STATUS_INTERNAL_SERVER_ERROR, err)
	}
	ts.Commit()

	return profile, nil
}

// JoinGroup - user join the group
func (s *GroupService) JoinGroup(userId, groupId, reason string) (profile *models.Group, tErr *TError) {
	isMember, err := s.isGroupMember(groupId, userId)
	if err != nil {
		return nil, NewTError(protos.StatusCode_STATUS_INTERNAL_SERVER_ERROR, err)
	}
	if isMember {
		return nil, NewTError(protos.StatusCode_STATUS_REPEAT_OPERATION, errmsg.ErrRepeatedlyJoinGroup)
	}

	newMember := models.NewMember(userId, "")
	newMember.GroupId = groupId

	if err := groupRepo.InsertMembers(&newMember); err != nil {
		return nil, NewTError(protos.StatusCode_STATUS_INTERNAL_SERVER_ERROR, err)
	}

	var condition = map[string]interface{}{
		"groupId": groupId,
	}
	profile, err = groupRepo.FindOneGroup(condition)
	if err != nil {
		return nil, NewTError(protos.StatusCode_STATUS_INTERNAL_SERVER_ERROR, err)
	}
	return profile, nil
}

// QuitGroup - user quit the group
func (s *GroupService) QuitGroup(userId, groupId string) (tErr *TError) {

	isMember, err := s.isGroupMember(groupId, userId)
	if err != nil {
		return NewTError(protos.StatusCode_STATUS_INTERNAL_SERVER_ERROR, err)
	}
	if !isMember {
		return NewTError(protos.StatusCode_STATUS_REPEAT_OPERATION, errmsg.ErrNotGroupMembers)
	}

	memberIds := []string{userId}
	count, err := groupRepo.RemoveMembers(groupId, memberIds, true)
	if err != nil || count == 0 { // current stage count should not be 0
		return NewTError(protos.StatusCode_STATUS_INTERNAL_SERVER_ERROR, err)
	}
	return
}

// AddGroupMember - add multiple members into group
func (s *GroupService) AddGroupMember(groupId string, memberIds []string) (profile *models.Group, tErr *TError) {
	condition := map[string]interface{}{
		"groupId": groupId,
	}
	groupProfile, err := groupRepo.FindOneGroup(condition)
	if err != nil {
		return nil, NewTError(protos.StatusCode_STATUS_INTERNAL_SERVER_ERROR, errmsg.ErrNotGroupMembers)
	}
	if groupProfile == nil {
		return nil, NewTError(protos.StatusCode_STATUS_BAD_REQUEST, errmsg.ErrGroupNotExists)
	}
	members := groupProfile.Members
	if isOutOfMemberLimit(len(members), len(memberIds)) {
		return nil, NewTError(protos.StatusCode_STATUS_BAD_REQUEST, errmsg.ErrGroupMemberReachedLimit)
	}

	newMembers := buildMembersObject(groupId, memberIds)
	if err = groupRepo.InsertMembers(newMembers...); err != nil {
		return nil, NewTError(protos.StatusCode_STATUS_INTERNAL_SERVER_ERROR, err)
	}
	profile, err = groupRepo.FindOneGroup(condition)
	if err != nil {
		return nil, NewTError(protos.StatusCode_STATUS_INTERNAL_SERVER_ERROR, err)
	}
	return profile, nil
}

// AddGroupMember - kick one member from a group
func (s *GroupService) KickGroupMember(userId, groupId, memberUserId string) (profile *models.Group, tErr *TError) {
	condition := map[string]interface{}{
		"groupId": groupId,
	}
	groupProfile, err := groupRepo.FindOneGroup(condition)
	if err != nil {
		return nil, NewTError(protos.StatusCode_STATUS_INTERNAL_SERVER_ERROR, err)
	}
	if groupProfile == nil {
		return nil, NewTError(protos.StatusCode_STATUS_BAD_REQUEST, errmsg.ErrGroupNotExists)
	}

	isMember, err := s.isGroupMember(groupId, userId)
	if err != nil {
		return nil, NewTError(protos.StatusCode_STATUS_INTERNAL_SERVER_ERROR, err)
	}
	if !isMember {
		return nil, NewTError(protos.StatusCode_STATUS_REPEAT_OPERATION, errmsg.ErrNotGroupMembers)
	}

	// check out currently user permission
	if isGroupManager(userId, groupProfile.OwnerUserId) {
		return nil, NewTError(protos.StatusCode_STATUS_REQUEST_FORBIDDEN, errmsg.ErrOperationForbidden)
	}

	_, err = groupRepo.RemoveMembers(groupId, []string{memberUserId}, true)
	if err != nil {
		return nil, NewTError(protos.StatusCode_STATUS_INTERNAL_SERVER_ERROR, err)
	}

	profile, err = groupRepo.FindOneGroup(condition)
	if err != nil {
		return nil, NewTError(protos.StatusCode_STATUS_INTERNAL_SERVER_ERROR, err)
	}
	return profile, nil
}

// UpdateGroupName - update group name by group manager
func (s *GroupService) UpdateGroupName(userId, groupId, newGroupName string) (profile *models.Group, tErr *TError) {
	// handle empty group name to default
	if newGroupName == "" {
		newGroupName = models.DefaultGroupName
	}

	condition := map[string]interface{}{
		"groupId": groupId,
	}
	groupProfile, err := groupRepo.FindOneGroup(condition)
	if err != nil {
		return nil, NewTError(protos.StatusCode_STATUS_INTERNAL_SERVER_ERROR, err)
	}
	if groupProfile == nil {
		return nil, NewTError(protos.StatusCode_STATUS_BAD_REQUEST, errmsg.ErrGroupNotExists)
	}

	// check out user permission
	if !isGroupManager(userId, groupProfile.OwnerUserId) {
		return nil, NewTError(protos.StatusCode_STATUS_REQUEST_FORBIDDEN, errmsg.ErrOperationForbidden)
	}

	updated := map[string]interface{}{
		"name": newGroupName,
	}
	profile, err = groupRepo.FindOneGroupAndUpdate(condition, updated)
	if err != nil {
		return nil, NewTError(protos.StatusCode_STATUS_INTERNAL_SERVER_ERROR, err)
	}
	return profile, nil
}

// UpdateGroupNotice - update group notice by group manager
func (s *GroupService) UpdateGroupNotice(userId, groupId, newNotice string) (profile *models.Group, tErr *TError) {
	condition := map[string]interface{}{
		"groupId": groupId,
	}
	groupProfile, err := groupRepo.FindOneGroup(condition)
	if err != nil {
		return nil, NewTError(protos.StatusCode_STATUS_INTERNAL_SERVER_ERROR, err)
	}
	if groupProfile == nil {
		return nil, NewTError(protos.StatusCode_STATUS_BAD_REQUEST, errmsg.ErrGroupNotExists)
	}
	// check out user permission
	if !isGroupManager(userId, groupProfile.OwnerUserId) {
		return nil, NewTError(protos.StatusCode_STATUS_REQUEST_FORBIDDEN, errmsg.ErrOperationForbidden)
	}

	updated := map[string]interface{}{
		"notice": newNotice,
	}
	profile, err = groupRepo.FindOneGroupAndUpdate(condition, updated)
	if err != nil {
		return nil, NewTError(protos.StatusCode_STATUS_INTERNAL_SERVER_ERROR, err)
	}
	return profile, nil
}

// UpdateMemberNickname - update current user's alias in group
func (s *GroupService) UpdateMemberNickname(userId, groupId, newAlias string) (profile *models.Group, tErr *TError) {
	condition := map[string]interface{}{
		"groupId": groupId,
		"userId":  userId,
	}
	updated := map[string]interface{}{
		"alias": newAlias,
	}
	memberProfile, err := groupRepo.FindOneMemberAndUpdate(condition, updated)
	if err != nil {
		return nil, NewTError(protos.StatusCode_STATUS_INTERNAL_SERVER_ERROR, err)
	}
	if memberProfile == nil {
		return nil, NewTError(protos.StatusCode_STATUS_BAD_REQUEST, errmsg.ErrNotGroupMembers)
	}

	searchGroup := map[string]interface{}{
		"groupId": groupId,
	}
	profile, err = groupRepo.FindOneGroup(searchGroup)
	if err != nil {
		return nil, NewTError(protos.StatusCode_STATUS_INTERNAL_SERVER_ERROR, err)
	}
	if profile == nil { // current stage profile would not be nil
		return nil, NewTError(protos.StatusCode_STATUS_INTERNAL_SERVER_ERROR, errmsg.ErrSystemUncheckException)
	}
	return profile, nil
}

// isGroupMember - check target user whether is a member of the group
func (s *GroupService) isGroupMember(groupId, userId string) (bool, error) {
	var condition = map[string]interface{}{
		"groupId": groupId,
		"userId":  userId,
	}
	memberProfile, err := groupRepo.FindOneMember(condition)
	if err != nil {
		return false, err
	}
	return memberProfile == nil, nil
}

// isGroupCountOverflow - check target user's group whether is reach limit count
func (s *GroupService) isGroupCountOverflow(userId string) (isOverflow bool, err error) {
	condition := map[string]interface{}{
		"ownerUserId": userId,
	}
	totalNum, err := groupRepo.CountGroup(condition)
	if err != nil {
		return false, err
	}
	if totalNum >= MaximumNumberOfGroups {
		return true, nil
	}
	return false, nil
}

func buildGroupProfile(userId string, groupName string, memberIds []string) *models.Group {
	var members = make([]models.Member, len(memberIds))
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
