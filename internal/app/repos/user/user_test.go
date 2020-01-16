package user

import (
	. "github.com/smartystreets/goconvey/convey"
	mysqlsrv "goim-pro/pkg/db/mysql"
	"testing"
)

var user1 = &User{
	Password: "1234567890",
	UserProfile: UserProfile{
		UserID:      2,
		Telephone:   "13631210010",
		Email:       "294001@qq.com",
		Username:    "TEST02",
		Nickname:    "TEST02",
		Description: "Never Settle",
		Birthday:    1578903121862,
	},
}

var user2 = &User{
	Password: "1234567890",
	UserProfile: UserProfile{
		UserID:      3,
		Telephone:   "13631210022",
		Email:       "294001@qq.com",
		Username:    "TEST02",
		Nickname:    "TEST02",
		Description: "Never Settle",
	},
}

func TestUser_IsTelephoneRegistered(t *testing.T) {
	mysqlDB := mysqlsrv.NewMysqlConnection()
	_ = mysqlDB.Connect()
	NewUserRepo(mysqlsrv.NewMysqlConnection().GetMysqlInstance())
	u := &User{}
	_ = u.Register(user2) // create a user
	Convey("Test_IsTelephoneRegistered", t, func() {
		Convey("Test_return_FALSE", func() {
			isExist, err := u.IsTelephoneRegistered("13631210033")
			So(err, ShouldBeNil)
			So(isExist, ShouldBeFalse)
		})
		Convey("Test_return_TRUE", func() {
			isExist, err := u.IsTelephoneRegistered("13631210022")
			So(err, ShouldBeNil)
			So(isExist, ShouldBeTrue)
		})
	})
	_ = u.RemoveUserByUserID(user1.UserID, true) // remove demo user
}

func TestUser_Register(t *testing.T) {
	mysqlDB := mysqlsrv.NewMysqlConnection()
	_ = mysqlDB.Connect()
	NewUserRepo(mysqlsrv.NewMysqlConnection().GetMysqlInstance())
	u := &User{}
	Convey("Test_Register", t, func() {
		Convey("Registration_successful", func() {
			err := u.Register(user1)
			So(err, ShouldBeNil)
		})
		Convey("Registration_fail_by_telephone_exist", func() {
			err := u.Register(user1)
			So(err, ShouldNotBeNil)
		})
	})

	_ = u.RemoveUserByUserID(user1.UserID, true)
}

func TestUser_RemoveUserByUserID(t *testing.T) {
	mysqlDB := mysqlsrv.NewMysqlConnection()
	_ = mysqlDB.Connect()
	NewUserRepo(mysqlsrv.NewMysqlConnection().GetMysqlInstance())

	Convey("TestUserRepo_RemoveUserByUserID", t, func() {
		Convey("testing_RemoveUserByUserID_success", func() {
			user := &User{
				Status: "ACTIVE",
				UserProfile: UserProfile{
					UserID: 2,
				},
			}
			err := user.RemoveUserByUserID(1, false)
			So(err, ShouldBeNil)
		})
	})
}
