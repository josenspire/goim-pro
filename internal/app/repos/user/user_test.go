package user

import (
	"goim-pro/config"
	mysqlsrv "goim-pro/pkg/db/mysql"
	"goim-pro/pkg/utils"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var user1 = &User{
	Password: "1234567890",
	UserProfile: UserProfile{
		UserID:      "2",
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
		UserID:      "3",
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
		Convey("Test_return_FALSE_with_exist_telephone", func() {
			isExist, err := u.IsTelephoneOrEmailRegistered("13631210033", "123@qq.com")
			So(err, ShouldBeNil)
			So(isExist, ShouldBeFalse)
		})
		Convey("Test_return_FALSE_with_exist_email", func() {
			isExist, err := u.IsTelephoneOrEmailRegistered("13631210044", "294001@qq.com")
			So(err, ShouldBeNil)
			So(isExist, ShouldBeTrue)
		})
		Convey("Test_return_TRUE", func() {
			isExist, err := u.IsTelephoneOrEmailRegistered("13631210022", "")
			So(err, ShouldBeNil)
			So(isExist, ShouldBeTrue)
		})
	})
	_ = u.RemoveUserByUserID(user2.UserID, true) // remove demo user
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
					UserID: "2",
				},
			}
			err := user.RemoveUserByUserID("1", false)
			So(err, ShouldBeNil)
		})
	})
}

func TestUser_Login(t *testing.T) {
	mysqlDB := mysqlsrv.NewMysqlConnection()
	_ = mysqlDB.Connect()
	NewUserRepo(mysqlsrv.NewMysqlConnection().GetMysqlInstance())

	u := &User{}
	_ = u.Register(user2) // create a user

	Convey("TestUserRepo_LoginByTelephone", t, func() {
		Convey("login_fail_with_incorrect_telephone_and_password", func() {
			_, err := u.LoginByTelephone("13631210022", "1234567890")
			So(err, ShouldNotBeNil)
			So(err, ShouldEqual, utils.ErrAccountOrPswInvalid)
		})
		Convey("login_success_then_return_userProfile", func() {
			enPassword, _ := crypto.AESEncrypt("1234567890", config.GetApiSecretKey())
			currUser, err := u.LoginByTelephone("13631210022", enPassword)
			So(err, ShouldBeNil)
			So(currUser.UserID, ShouldEqual, "3")
			So(currUser.Telephone, ShouldEqual, "13631210022")
		})
	})
	_ = u.RemoveUserByUserID(user2.UserID, true) // remove demo user
}

func TestUser_LoginByEmail(t *testing.T) {
	mysqlDB := mysqlsrv.NewMysqlConnection()
	_ = mysqlDB.Connect()
	NewUserRepo(mysqlsrv.NewMysqlConnection().GetMysqlInstance())

	u := &User{}
	_ = u.Register(user2) // create a user
	Convey("TestUserRepo_LoginByEmail", t, func() {
		Convey("login_fail_with_incorrect_email_and_password", func() {
			_, err := u.LoginByEmail("294001@qq.com", "1234567890")
			So(err, ShouldNotBeNil)
			So(err, ShouldEqual, utils.ErrAccountOrPswInvalid)
		})
		Convey("login_success_then_return_userProfile", func() {
			enPassword, _ := crypto.AESEncrypt("1234567890", config.GetApiSecretKey())
			currUser, err := u.LoginByEmail("294001@qq.com", enPassword)
			So(err, ShouldBeNil)
			So(currUser.UserID, ShouldEqual, "3")
			So(currUser.Email, ShouldEqual, "294001@qq.com")
		})
	})
	_ = u.RemoveUserByUserID(user2.UserID, true) // remove demo user
}
