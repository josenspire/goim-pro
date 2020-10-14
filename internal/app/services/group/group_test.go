package groupsrv

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
	consts "goim-pro/internal/app/constants"
	"goim-pro/internal/app/models"
	mysqlsrv "goim-pro/pkg/db/mysql"
	errmsg "goim-pro/pkg/errors"
	"testing"
)

func TestGroupService_CreateGroup(t *testing.T) {
	_ = mysqlsrv.NewMysql()

	var members = []models.Member{
		models.Member{
			UserId:  "TEST001",
			Alias:   "",
			Role:    "1",
			Status:  models.DefaultGroupMemberStatus,
			GroupId: "TEST_GROUP_001",
		},
		models.Member{
			UserId:  "TEST002",
			Alias:   "",
			Role:    "1",
			Status:  models.DefaultGroupMemberStatus,
			GroupId: "TEST_GROUP_001",
		},
	}
	var groupProfile1 = &models.Group{
		GroupId:     "TEST_GROUP_001",
		OwnerUserId: "TEST001",
		Name:        "GROUP_NAME_01",
		Members:     members,
	}
	mg := &MockGroup{}
	mg.On("CountGroup", map[string]interface{}{"ownerUserId": "TEST001"}).Return(3, nil)
	mg.On("CountGroup", map[string]interface{}{"ownerUserId": "TEST003"}).Return(consts.MaximumNumberOfGroups, nil)

	mg.On("CreateGroup", mock.Anything).Return(groupProfile1, nil)

	mg.On("FindOneGroup", map[string]interface{}{"groupId": "TEST_GROUP_001"}).Return(groupProfile1, nil)

	gs := New()
	groupRepo = mg
	Convey("Testing_CreateGroup", t, func() {
		Convey("should_create_group_successful_when_given_correct_members_then_return_group_profile", func() {
			profile, err := gs.CreateGroup("TEST001", "GROUP_NAME_01", []string{"TEST002"})

			So(err, ShouldBeNil)
			So(profile.GroupId, ShouldEqual, "TEST_GROUP_001")
		})
		Convey("should_create_group_failed_when_group_only_include_current_user_then_return_error", func() {
			profile, err := gs.CreateGroup("TEST002", "GROUP_NAME_01", []string{"TEST002"})

			So(profile, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err.Detail, ShouldEqual, errmsg.ErrIllegalOperation.Error())
		})
		Convey("should_create_group_failed_when_current_user's_group_reached_limit_then_return_error", func() {
			profile, err := gs.CreateGroup("TEST003", "GROUP_NAME_01", []string{"TEST002"})

			So(profile, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err.Detail, ShouldEqual, errmsg.ErrGroupReachedLimit.Error())
		})
	})
}

func TestGroupService_AddGroupMember(t *testing.T) {
}

func TestGroupService_JoinGroup(t *testing.T) {
}

func TestGroupService_KickGroupMember(t *testing.T) {
}

func TestGroupService_QuitGroup(t *testing.T) {
}

func TestGroupService_UpdateGroupName(t *testing.T) {
}

func TestGroupService_UpdateGroupNotice(t *testing.T) {
}

func TestGroupService_UpdateMemberNickname(t *testing.T) {
}

func TestGroupService_isGroupCountOverflow(t *testing.T) {
}

func TestGroupService_isGroupMember(t *testing.T) {
}
