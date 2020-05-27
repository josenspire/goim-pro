package group

import (
	"goim-pro/internal/app/models"
	mysqlsrv "goim-pro/pkg/db/mysql"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGroupImpl_CreateGroup(t *testing.T) {
	mysqlDB := mysqlsrv.NewMysql()
	_ = mysqlDB.Connect()
	s := NewGroupRepo(mysqlsrv.NewMysql())

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

func TestGroupImpl_FindOneGroup(t *testing.T) {
	mysqlDB := mysqlsrv.NewMysql()
	_ = mysqlDB.Connect()
	NewGroupRepo(mysqlsrv.NewMysql())

	newMember1 := models.NewMember("01E59Z8HMG8SK8C65XV42M33QP", "JAMES_TEST_001")
	newMember2 := models.NewMember("01E59ZNYB8KDNW0W3NHGDZDD6V", "JAMES_TEST_002")
	members := []models.Member{
		newMember1,
		newMember2,
	}
	groupProfile := models.NewGroup("TEST_GROUP_001", "TEST005", "TEST_GROUP_001", members)

	s := &GroupImpl{}
	_, err := s.CreateGroup(groupProfile)
	if err != nil {
		t.FailNow()
	}
	Convey("Test_FindOneGroup", t, func() {
		Convey("should_find_one_group_with_nil_err", func() {
			condition := map[string]interface{}{
				"groupId":     "TEST_GROUP_001",
				"ownerUserId": "TEST005",
			}
			profile, err := s.FindOneGroup(condition)
			ShouldBeNil(err)
			So(profile.OwnerUserId, ShouldEqual, "TEST005")
			So(len(profile.Members), ShouldEqual, 2)
			So(profile.Members[0].User.Nickname, ShouldEqual, "JAMES01")
		})
		Convey("should_not_find_the_group_when_given_error_groupId_or_ownerUserId_then_return_nil", func() {
			condition := map[string]interface{}{
				"groupId":     "TEST_GROUP_000001",
				"ownerUserId": "TEST0005",
			}
			profile, err := s.FindOneGroup(condition)
			ShouldBeNil(err)
			ShouldBeNil(profile)
		})
	})

	_ = s.RemoveGroupByGroupId("TEST_GROUP_001", true)
}

func TestGroupImpl_CountGroup(t *testing.T) {
	mysqlDB := mysqlsrv.NewMysql()
	_ = mysqlDB.Connect()
	db := mysqlsrv.NewMysql()
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
				"ownerUserId": "TEST005",
			}
			count, err := s.CountGroup(condition)

			ShouldBeNil(err)
			So(count, ShouldEqual, 2)
		})
	})
	_ = s.RemoveGroupByGroupId("TEST_GROUP_01", true)
	_ = s.RemoveGroupByGroupId("TEST_GROUP_02", true)
}

func TestGroupImpl_FindOneGroupAndUpdate(t *testing.T) {
	mysqlDB := mysqlsrv.NewMysql()
	_ = mysqlDB.Connect()
	db := mysqlsrv.NewMysql()
	NewGroupRepo(db)

	newMember1 := models.NewMember("TEST001", "JAMES_TEST_001")
	newMember2 := models.NewMember("TEST002", "JAMES_TEST_002")
	members := []models.Member{
		newMember1,
		newMember2,
	}

	groupProfile := models.NewGroup("TEST_GROUP_01", "TEST005", "TEST_GROUP_001", members)
	groupProfile.Notice = "Never Settle!"
	s := &GroupImpl{}
	if _, err := s.CreateGroup(groupProfile); err != nil {
		t.FailNow()
	}
	Convey("Test_FindOneGroupAndUpdate", t, func() {
		Convey("should_find_and_update_group_name_success_then_return_new_group_profile_with_nil_error", func() {
			condition := map[string]interface{}{
				"groupId":     "TEST_GROUP_01",
				"ownerUserId": "TEST005",
			}
			updated := map[string]interface{}{
				"name": "NEW_GROUP_NAME_01",
			}
			newProfile, err := s.FindOneGroupAndUpdate(condition, updated)

			ShouldBeNil(err)
			So(newProfile.Name, ShouldEqual, "NEW_GROUP_NAME_01")
			So(newProfile.Notice, ShouldEqual, "Never Settle!")
		})
		Convey("should_not_find_the_group_when_given_incorrect_groupId_or_ownerUserId_then_return_nil_result_with_nil_err", func() {
			condition := map[string]interface{}{
				"groupId":     "TEST_GROUP_02",
				"ownerUserId": "TEST005",
			}
			updated := map[string]interface{}{
				"name": "NEW_GROUP_NAME_01",
			}
			newProfile, err := s.FindOneGroupAndUpdate(condition, updated)

			ShouldBeNil(err)
			ShouldBeNil(newProfile)
		})
	})
	_ = s.RemoveGroupByGroupId("TEST_GROUP_01", true)
}

func TestGroupImpl_FindOneMember(t *testing.T) {
	mysqlDB := mysqlsrv.NewMysql()
	_ = mysqlDB.Connect()
	NewGroupRepo(mysqlsrv.NewMysql())

	newMember1 := models.NewMember("TEST001", "JAMES_TEST_001")
	newMember2 := models.NewMember("TEST002", "JAMES_TEST_002")
	members := []models.Member{
		newMember1,
		newMember2,
	}
	groupProfile := models.NewGroup("TEST_GROUP_001", "TEST005", "TEST_GROUP_001", members)

	s := &GroupImpl{}
	_, err := s.CreateGroup(groupProfile)
	if err != nil {
		t.FailNow()
	}

	Convey("Test_FindOneGroupMember", t, func() {
		Convey("should_return_one_member_profile_with_nil_error", func() {
			condition := map[string]interface{}{
				"groupId": "TEST_GROUP_001",
				"userId":  "TEST002",
			}
			memberProfile, err := s.FindOneMember(condition)
			ShouldBeNil(err)
			So(memberProfile.GroupId, ShouldEqual, "TEST_GROUP_001")
			So(memberProfile.Alias, ShouldEqual, "JAMES_TEST_002")
		})

		Convey("should_return_nil_result_with_nil_error_when_given_not_exists_groupId_or_userId", func() {
			condition := map[string]interface{}{
				"groupId": "TEST_GROUP_002",
				"userId":  "TEST002",
			}
			memberProfile, err := s.FindOneMember(condition)
			ShouldBeNil(err)
			ShouldBeNil(memberProfile)
		})
	})

	_ = s.RemoveGroupByGroupId("TEST_GROUP_001", true)
}

func TestGroupImpl_InsertMembers(t *testing.T) {
	mysqlDB := mysqlsrv.NewMysql()
	_ = mysqlDB.Connect()
	NewGroupRepo(mysqlsrv.NewMysql())

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
			ShouldBeNil(err)

			condition := map[string]interface{}{
				"groupId": group.GroupId,
				"userId":  newMember3.UserId,
			}
			memberProfile, err := s.FindOneMember(condition)
			ShouldBeNil(err)
			So(memberProfile.UserId, ShouldEqual, newMember3.UserId)
			So(memberProfile.GroupId, ShouldEqual, group.GroupId)
		})
	})

	_ = s.RemoveGroupByGroupId("TEST_GROUP_001", true)
}

func TestGroupImpl_RemoveMembers(t *testing.T) {
	mysqlDB := mysqlsrv.NewMysql()
	_ = mysqlDB.Connect()
	NewGroupRepo(mysqlsrv.NewMysql())

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
			count, err := s.RemoveMembers(group.GroupId, memberIds, true)
			ShouldBeNil(err)
			So(count, ShouldEqual, 2)

			condition := map[string]interface{}{
				"groupId": "TEST_GROUP_001",
			}
			_group, err := s.FindOneGroup(condition)
			ShouldBeNil(err)
			So(_group.OwnerUserId, ShouldEqual, "TEST005")
		})
	})

	_ = s.RemoveGroupByGroupId("TEST_GROUP_001", true)
}

func TestGroupImpl_FindOneMemberAndUpdate(t *testing.T) {
	mysqlDB := mysqlsrv.NewMysql()
	_ = mysqlDB.Connect()
	db := mysqlsrv.NewMysql()
	NewGroupRepo(db)

	newMember1 := models.NewMember("TEST001", "JAMES_TEST_001")
	newMember2 := models.NewMember("TEST002", "JAMES_TEST_002")
	members := []models.Member{
		newMember1,
		newMember2,
	}

	groupProfile := models.NewGroup("TEST_GROUP_01", "TEST005", "TEST_GROUP_001", members)
	groupProfile.Notice = "Never Settle!"
	s := &GroupImpl{}
	if _, err := s.CreateGroup(groupProfile); err != nil {
		t.FailNow()
	}
	Convey("Test_FindOneGroupAndUpdate", t, func() {
		Convey("should_find_and_update_member_alias_success_then_return_new_member_profile_with_nil_error", func() {
			condition := map[string]interface{}{
				"groupId": "TEST_GROUP_01",
				"userId":  "TEST002",
			}
			updated := map[string]interface{}{
				"alias": "NEW_ALIAS",
			}
			newProfile, err := s.FindOneMemberAndUpdate(condition, updated)

			ShouldBeNil(err)
			So(newProfile.UserId, ShouldEqual, "TEST002")
			So(newProfile.Alias, ShouldEqual, "NEW_ALIAS")
		})
		Convey("should_not_find_the_member_when_given_incorrect_groupId_or_userId_then_return_nil_result_with_nil_err", func() {
			condition := map[string]interface{}{
				"groupId": "TEST_GROUP_01",
				"userId":  "TEST003",
			}
			updated := map[string]interface{}{
				"alias": "NEW_ALIAS",
			}
			newProfile, err := s.FindOneMemberAndUpdate(condition, updated)

			ShouldBeNil(err)
			ShouldBeNil(newProfile)
		})
	})
	_ = s.RemoveGroupByGroupId("TEST_GROUP_01", true)
}
