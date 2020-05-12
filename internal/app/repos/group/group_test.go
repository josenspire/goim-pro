package group

import (
	. "github.com/smartystreets/goconvey/convey"
	"goim-pro/internal/app/models"
	mysqlsrv "goim-pro/pkg/db/mysql"
	"testing"
)

func TestGroupImpl_CreateGroup(t *testing.T) {
	mysqlDB := mysqlsrv.NewMysqlConnection()
	_ = mysqlDB.Connect()
	s := NewGroupRepo(mysqlsrv.NewMysqlConnection().GetMysqlInstance())

	newMember1 := models.NewMember("TEST001", "JAMES_TEST_001")
	newMember2 := models.NewMember("TEST002", "JAMES_TEST_002")
	members := []models.Member{
		newMember1, newMember2,
	}
	groupProfile := models.NewGroup("TEST005", "TEST_GROUP_001", members)

	var profile = &models.Group{}
	var err error
	Convey("Test_InsertGroup", t, func() {
		Convey("should_create_group_with_members_then_return_group_profile", func() {
			profile, err = s.CreateGroup(groupProfile)
			ShouldBeNil(err)
			So(profile.Name, ShouldEqual, "TEST_GROUP_001")
		})
	})

	_ = s.RemoveGroupByGroupId(profile.GroupId, true)
}

func TestGroupImpl_InsertMembers(t *testing.T) {
	mysqlDB := mysqlsrv.NewMysqlConnection()
	_ = mysqlDB.Connect()
	NewGroupRepo(mysqlsrv.NewMysqlConnection().GetMysqlInstance())

	newMember1 := models.NewMember("TEST001", "JAMES_TEST_001")
	newMember2 := models.NewMember("TEST002", "JAMES_TEST_002")
	members := []models.Member{
		newMember1,
		newMember2,
	}
	groupProfile := models.NewGroup("TEST005", "TEST_GROUP_001", members)

	var group *models.Group
	var err error
	s := &GroupImpl{}
	group, err = s.CreateGroup(groupProfile)
	if err != nil {
		t.FailNow()
	}

	Convey("Test_InsertMembers", t, func() {
		Convey("should_insert_success_and_return_nil_of_error", func() {
			newMember3 := models.NewMember("TEST003", "JAMES_TEST_003")
			newMember3.GroupId = group.GroupId

			err = s.InsertMembers(&newMember3)
			if err != nil {
				t.FailNow()
			}

			ShouldBeNil(err)
			ShouldBeTrue(isExist)
		})
	})

	_ = ct.RemoveContactsByIds("TEST001", "TEST002")
}

func TestGroupImpl_FindOneGroup(t *testing.T) {

}

func TestGroupImpl_FindOneGroupMember(t *testing.T) {

}

func TestGroupImpl_InsertMembers1(t *testing.T) {

}

func TestGroupImpl_RemoveGroupByGroupId(t *testing.T) {

}
