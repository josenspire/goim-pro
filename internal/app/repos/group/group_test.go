package group

import (
	. "github.com/smartystreets/goconvey/convey"
	"goim-pro/internal/app/models"
	mysqlsrv "goim-pro/pkg/db/mysql"
	"testing"
)

func TestGroupImpl_InsertGroup(t *testing.T) {
	mysqlDB := mysqlsrv.NewMysqlConnection()
	_ = mysqlDB.Connect()
	s := NewGroupRepo(mysqlsrv.NewMysqlConnection().GetMysqlInstance())

	newMember1 := models.NewMember("TEST001", "JAMES_TEST_001")
	newMember2 := models.NewMember("TEST002", "JAMES_TEST_002")

	groupProfile := models.NewGroup("TEST005", "TEST_GROUP_001")
	groupProfile.Members = []models.Member{
		newMember1, newMember2,
	}

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
	//mysqlDB := mysqlsrv.NewMysqlConnection()
	//_ = mysqlDB.Connect()
	//NewGroupRepo(mysqlsrv.NewMysqlConnection().GetMysqlInstance())
	//
	//groupProfile := models.NewGroup("TEST001", "")
	//newMember1 := models.NewMember("TEST001", "JAMES_TEST_001")
	//
	//var group *models.Group
	//var err error
	//s := &GroupImpl{}
	//if group, err = s.CreateGroup(groupProfile); err != nil {
	//	t.FailNow()
	//}
	//err = s.InsertMembers(newMember1)
	//if err != nil {
	//	t.FailNow()
	//}
	//
	//Convey("Test_InsertMembers", t, func() {
	//	Convey("should_return_true_result", func() {
	//		ShouldBeNil(err)
	//		ShouldBeTrue(isExist)
	//	})
	//})
	//
	//_ = ct.RemoveContactsByIds("TEST001", "TEST002")
}
