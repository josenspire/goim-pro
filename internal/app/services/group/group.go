package group

import (
	"context"
	"github.com/jinzhu/gorm"
	protos "goim-pro/api/protos/salty"
	. "goim-pro/internal/app/constants"
	"goim-pro/internal/app/repos"
	. "goim-pro/internal/app/repos/group"
	. "goim-pro/internal/app/repos/user"
	"goim-pro/internal/app/services/converters"
	mysqlsrv "goim-pro/pkg/db/mysql"
	redsrv "goim-pro/pkg/db/redis"
	"goim-pro/pkg/logs"
	"goim-pro/pkg/utils"
	"net/http"
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

	repoServer := repos.New(mysqlDB)
	return &groupService{
		userRepo:  repoServer.UserRepo,
		groupRepo: repoServer.GroupRepo,
	}
}

// create a group
func (gs *groupService) CreateGroup(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
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

	totalGroup, err := isGroupCountOverflow(gs, userId)
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

	groupProfile := NewGroup(userId, groupName)

	ts := mysqlDB.Begin()
	groupProfile, err = gs.groupRepo.InsertOne(groupProfile)
	if err != nil {
		ts.Callback()

		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		return
	}
	newMembers := buildMembersObject(memberIds)
	members, err := gs.groupRepo.InsertMembers(newMembers...)
	// TODO:
	logger.Info(members)
	if err != nil {
		ts.Callback()

		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		return
	}

	registerResp := &protos.CreateGroupResp{
		Profile: converters.ConvertEntity2ProtoForGroupProfile(groupProfile),
	}
	resp.Data, err = utils.MarshalMessageToAny(registerResp)
	if err != nil {
		logger.Errorf("register response marshal message error: %s", err.Error())
	}
	resp.Message = "user registration successful"
	return
}

func (gs *groupService) JoinGroup(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	panic("implement me")
}

func (gs *groupService) QuitGroup(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	panic("implement me")
}

func (gs *groupService) AddGroupMember(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	panic("implement me")
}

func (gs *groupService) KickGroupMember(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	panic("implement me")
}

func (gs *groupService) UpdateGroupName(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	panic("implement me")
}

func (gs *groupService) UpdateGroupNotice(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	panic("implement me")
}

func (gs *groupService) UpdateMemberNickname(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	panic("implement me")
}

func isGroupCountOverflow(gs *groupService, userId string) (totalNum int, err error) {
	return
}

func buildMembersObject(memberIds []string) (members []*Member) {
	members = make([]*Member, len(memberIds)-1)
	for i, userId := range memberIds {
		members[i] = NewMember(userId, "")
	}
	return
}
