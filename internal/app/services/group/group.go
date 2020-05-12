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
	redsrv "goim-pro/pkg/db/redis"
	"goim-pro/pkg/http"
	"goim-pro/pkg/logs"
	"goim-pro/pkg/utils"
	"strings"
)

var (
	logger  = logs.GetLogger("INFO")
	myRedis *redsrv.BaseClient
	mysqlDB *gorm.DB
)

type groupService struct {
	userRepo  IUserRepo
	groupRepo IGroupRepo
}

func New() protos.GroupServiceServer {
	myRedis = redsrv.NewRedisConnection().GetRedisClient()
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
		resp.Message = utils.ErrIllegalOperation.Error()
		return
	}

	totalGroup, err := isGroupCountOverflow(s, userId)
	if err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		return
	}
	if totalGroup > MaxGroupCount {
		resp.Code = http.StatusBadRequest
		resp.Message = utils.ErrGroupReachedLimit.Error()
		return
	}

	groupProfile := models.NewGroup(userId, groupName, buildMembersObject(memberIds))
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
		resp.Message = utils.ErrIllegalOperation.Error()
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
		resp.Message = "user has joined the group"
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

func (s *groupService) QuitGroup(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	panic("implement me")
}

func (s *groupService) AddGroupMember(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	panic("implement me")
}

func (s *groupService) KickGroupMember(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	panic("implement me")
}

func (s *groupService) UpdateGroupName(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	panic("implement me")
}

func (s *groupService) UpdateGroupNotice(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	panic("implement me")
}

func (s *groupService) UpdateMemberNickname(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	panic("implement me")
}

func (s *groupService) isGroupMember(groupId, userId string) (bool, error) {
	var condition = map[string]interface{}{
		"groupId": groupId,
		"userId":  userId,
	}
	_, err := s.groupRepo.FindOneGroupMember(condition)
	if err == utils.ErrGroupNotExists {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func isGroupCountOverflow(gs *groupService, userId string) (totalNum int, err error) {
	return
}

func buildMembersObject(memberIds []string) (members []models.Member) {
	members = make([]models.Member, len(memberIds)-1)
	for i, userId := range memberIds {
		members[i] = models.NewMember(userId, "")
	}
	return
}
