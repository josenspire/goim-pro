package group

import (
	"goim-pro/internal/app/models"
	mysqlsrv "goim-pro/pkg/db/mysql"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
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
	groupProfile := models.NewGroup("TEST_GROUP_001", "TEST005", "TEST_GROUP_001", members)

	var profile = &models.Group{}
	var err error
	Convey("Test_InsertGroup", t, func() {
		Convey("should_create_group_with_members_then_return_group_profile", func() {
			profile, err = s.CreateGroup(groupProfile)
			ShouldBeNil(err)
			So(profile.Name, ShouldEqual, "TEST_GROUP_001")
		})
	})

	_ = s.RemoveGroupByGroupId("TEST_GROUP_001", true)
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
	groupProfile := models.NewGroup("TEST_GROUP_001", "TEST005", "TEST_GROUP_001", members)

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

			condition := map[string]interface{}{
				"groupId": group.GroupId,
				"userId":  newMember3.UserId,
			}
			memberProfile, err := s.FindOneGroupMember(condition)
			ShouldBeNil(err)
			So(memberProfile.UserId, ShouldEqual, newMember3.UserId)
			So(memberProfile.GroupId, ShouldEqual, group.GroupId)
		})
	})

	_ = s.RemoveGroupByGroupId("TEST_GROUP_001", true)
}

func TestGroupImpl_FindOneGroup(t *testing.T) {

}

func TestGroupImpl_FindOneGroupMember(t *testing.T) {

}

func TestGroupImpl_InsertMembers1(t *testing.T) {

}

func TestGroupImpl_RemoveGroupByGroupId(t *testing.T) {

}

func TestGroupImpl_RemoveGroupMembers(t *testing.T) {
	mysqlDB := mysqlsrv.NewMysqlConnection()
	_ = mysqlDB.Connect()
	NewGroupRepo(mysqlsrv.NewMysqlConnection().GetMysqlInstance())

	newMember1 := models.NewMember("TEST001", "JAMES_TEST_001")
	newMember2 := models.NewMember("TEST002", "JAMES_TEST_002")
	members := []models.Member{
		newMember1,
		newMember2,
	}
	groupProfile := models.NewGroup("TEST_GROUP_001", "TEST005", "TEST_GROUP_001", members)

	var group *models.Group
	var err error
	s := &GroupImpl{}
	group, err = s.CreateGroup(groupProfile)
	if err != nil {
		t.FailNow()
	}

	Convey("Test_RemoveGroupMembers", t, func() {
		Convey("should_return_remove_count_and_nil_error", func() {
			memberIds := []string{
				newMember1.UserId,
				newMember2.UserId,
			}
			count, err := s.RemoveGroupMembers(group.GroupId, memberIds, true)

			ShouldBeNil(err)
			So(count, ShouldEqual, 2)
		})
	})

	_ = s.RemoveGroupByGroupId("TEST_GROUP_001", true)
}

func TestGroupImpl_CountGroup(t *testing.T) {
	mysqlDB := mysqlsrv.NewMysqlConnection()
	_ = mysqlDB.Connect()
	db := mysqlsrv.NewMysqlConnection().GetMysqlInstance()
	NewGroupRepo(db)

	newMember1 := models.NewMember("TEST001", "JAMES_TEST_001")
	newMember2 := models.NewMember("TEST002", "JAMES_TEST_002")
	members := []models.Member{
		newMember1,
		newMember2,
	}

	groupProfile1 := models.NewGroup("TEST_GROUP_01", "TEST005", "TEST_GROUP_001", members)
	groupProfile2 := models.NewGroup("TEST_GROUP_02", "TEST005", "TEST_GROUP_002", members)
	var err error
	s := &GroupImpl{}

	ts := db.Begin()
	_, err = s.CreateGroup(groupProfile1)
	if err != nil {
		ts.Rollback()
		t.FailNow()
	}
	_, err = s.CreateGroup(groupProfile2)
	if err != nil {
		ts.Rollback()
		t.FailNow()
	}
	ts.Commit()

	Convey("Test_CountGroup", t, func() {
		Convey("should_return_2_with_nil_error", func() {
			condition := map[string]interface{}{
				"userId": "TEST005",
			}
			count, err := s.CountGroup(condition)

			ShouldBeNil(err)
			So(count, ShouldEqual, 2)
		})
	})

	_ = s.RemoveGroupByGroupId("TEST_GROUP_001", true)
	_ = s.RemoveGroupByGroupId("TEST_GROUP_002", true)
}
