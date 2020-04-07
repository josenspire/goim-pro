package user

import (
	"goim-pro/config"
	"goim-pro/internal/app/models"
	mysqlsrv "goim-pro/pkg/db/mysql"
	"goim-pro/pkg/utils"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var user1 = &User{
	Password: "1234567890",
	UserProfile: models.UserProfile{
		UserId:      "2",
		Telephone:   "13631210010",
		Email:       "294001@qq.com",
		Nickname:    "TEST02",
		Description: "Never Settle",
		Birthday:    1578903121862,
	},
}

var user2 = &User{
	Password: "1234567890",
	UserProfile: models.UserProfile{
		UserId:      "3",
		Telephone:   "13631210022",
		Email:       "294001@qq.com",
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
	_ = u.RemoveUserByUserId(user2.UserId, true) // remove demo user
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

	_ = u.RemoveUserByUserId(user1.UserId, true)
}

func TestUser_RemoveUserByUserId(t *testing.T) {
	mysqlDB := mysqlsrv.NewMysqlConnection()
	_ = mysqlDB.Connect()
	NewUserRepo(mysqlsrv.NewMysqlConnection().GetMysqlInstance())

	Convey("TestUserRepo_RemoveUserByUserID", t, func() {
		Convey("testing_RemoveUserByUserID_success", func() {
			user := &User{
				Status: "ACTIVE",
				UserProfile: models.UserProfile{
					UserId: "2",
				},
			}
			err := user.RemoveUserByUserId("1", false)
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
			_, err := u.QueryByTelephoneAndPassword("13631210022", "1234567890")
			So(err, ShouldNotBeNil)
			So(err, ShouldEqual, utils.ErrAccountOrPwdInvalid)
		})
		Convey("login_success_then_return_userProfile", func() {
			enPassword, _ := crypto.AESEncrypt("1234567890", config.GetApiSecretKey())
			currUser, err := u.QueryByTelephoneAndPassword("13631210022", enPassword)
			So(err, ShouldBeNil)
			So(currUser.UserId, ShouldEqual, "3")
			So(currUser.Telephone, ShouldEqual, "13631210022")
		})
	})
	_ = u.RemoveUserByUserId(user2.UserId, true) // remove demo user
}

func TestUser_LoginByEmail(t *testing.T) {
	mysqlDB := mysqlsrv.NewMysqlConnection()
	_ = mysqlDB.Connect()
	NewUserRepo(mysqlsrv.NewMysqlConnection().GetMysqlInstance())

	u := &User{}
	_ = u.Register(user2) // create a user
	Convey("TestUserRepo_LoginByEmail", t, func() {
		Convey("login_fail_with_incorrect_email_and_password", func() {
			_, err := u.QueryByEmailAndPassword("294001@qq.com", "1234567890")
			So(err, ShouldNotBeNil)
			So(err, ShouldEqual, utils.ErrAccountOrPwdInvalid)
		})
		Convey("login_success_then_return_userProfile", func() {
			enPassword, _ := crypto.AESEncrypt("1234567890", config.GetApiSecretKey())
			currUser, err := u.QueryByEmailAndPassword("294001@qq.com", enPassword)
			So(err, ShouldBeNil)
			So(currUser.UserId, ShouldEqual, "3")
			So(currUser.Email, ShouldEqual, "294001@qq.com")
		})
	})
	_ = u.RemoveUserByUserId(user2.UserId, true) // remove demo user
}

func TestUser_ResetPasswordByTelephone(t *testing.T) {
	mysqlDB := mysqlsrv.NewMysqlConnection()
	_ = mysqlDB.Connect()
	NewUserRepo(mysqlsrv.NewMysqlConnection().GetMysqlInstance())

	u := &User{}
	_ = u.Register(user1) // create a user
	Convey("TestUserRepo_ResetPasswordByTelephone", t, func() {
		Convey("update_successful_by_telephone", func() {
			var telephone string = "13631210010"
			var password string = "111111111"
			newPassword, _ := crypto.AESEncrypt(password, config.GetApiSecretKey())

			err := u.ResetPasswordByTelephone(telephone, newPassword)
			So(err, ShouldBeNil)

			user, err := u.QueryByTelephoneAndPassword(telephone, newPassword)
			So(err, ShouldBeNil)
			So(user.Telephone, ShouldEqual, telephone)
		})
	})
	_ = u.RemoveUserByUserId(user1.UserId, true) // remove demo user
}

func TestUser_ResetPasswordByEmail(t *testing.T) {
	mysqlDB := mysqlsrv.NewMysqlConnection()
	_ = mysqlDB.Connect()
	NewUserRepo(mysqlsrv.NewMysqlConnection().GetMysqlInstance())

	u := &User{}
	_ = u.Register(user1) // create a user
	Convey("TestUserRepo_ResetPasswordByEmail", t, func() {
		Convey("update_successful_by_email", func() {
			var email string = "294001@qq.com"
			var password string = "2222222222"
			newPassword, _ := crypto.AESEncrypt(password, config.GetApiSecretKey())

			err := u.ResetPasswordByEmail(email, newPassword)
			So(err, ShouldBeNil)

			user, err := u.QueryByEmailAndPassword(email, newPassword)
			So(err, ShouldBeNil)
			So(user.Email, ShouldEqual, email)
		})
	})
	_ = u.RemoveUserByUserId(user1.UserId, true) // remove demo user
}

func TestUser_GetUserByUserId(t *testing.T) {
	mysqlDB := mysqlsrv.NewMysqlConnection()
	_ = mysqlDB.Connect()
	NewUserRepo(mysqlsrv.NewMysqlConnection().GetMysqlInstance())

	u := &User{}
	_ = u.Register(user1) // create a user
	Convey("TestUserRepo_GetUserByUserId", t, func() {
		Convey("get_user_info_success_then_return", func() {
			var userId = "2"

			user, err := u.FindByUserId(userId)
			So(err, ShouldBeNil)
			So(user.Email, ShouldEqual, "294001@qq.com")
		})
		Convey("get_user_info_failed_by_invalid_userId", func() {
			var userId = "3"

			_, err := u.FindByUserId(userId)
			So(err, ShouldNotBeNil)
			So(err, ShouldEqual, utils.ErrInvalidUserId)
		})
	})
	_ = u.RemoveUserByUserId(user1.UserId, true) // remove demo user
}

func TestUser_FindOneUser(t *testing.T) {
	mysqlDB := mysqlsrv.NewMysqlConnection()
	_ = mysqlDB.Connect()
	NewUserRepo(mysqlsrv.NewMysqlConnection().GetMysqlInstance())

	u := &User{}
	_ = u.Register(user1) // create a user
	Convey("TestUserRepo_FindOneUser", t, func() {
		Convey("get_user_info_success_by_telephone_then_return", func() {
			userCriteria := &User{}
			userCriteria.Telephone = "13631210010"

			user, err := u.FindOneUser(userCriteria)
			So(err, ShouldBeNil)
			So(user.Email, ShouldEqual, "294001@qq.com")
		})
		Convey("get_user_info_success_by_telephone_and_email_then_return", func() {
			userCriteria := &User{}
			userCriteria.Email = "294001@qq.com"
			userCriteria.Telephone = "13631210010"

			user, err := u.FindOneUser(userCriteria)
			So(err, ShouldBeNil)
			So(user.Email, ShouldEqual, "294001@qq.com")
		})
		Convey("get_user_info_failed_by_invalid_email_then_return_error", func() {
			userCriteria := &User{}
			userCriteria.Email = "294001@qq.com"
			userCriteria.Telephone = "13631210015"

			_, err := u.FindOneUser(userCriteria)
			So(err, ShouldNotBeNil)
			So(err, ShouldEqual, utils.ErrUserNotExists)
		})
	})
	_ = u.RemoveUserByUserId(user1.UserId, true) // remove demo user
}

func TestUser_FindOneAndUpdateProfile(t *testing.T) {
	mysqlDB := mysqlsrv.NewMysqlConnection()
	_ = mysqlDB.Connect()
	NewUserRepo(mysqlsrv.NewMysqlConnection().GetMysqlInstance())

	u := &User{}
	_ = u.Register(user1) // create a user
	Convey("TestUserRepo_FindOneAndUpdateProfile", t, func() {
		Convey("find_user_by_userId_and_update_profile_success", func() {
			userCriteria := &User{}
			userCriteria.UserId = "2"

			newProfile := models.UserProfile{
				UserId:      "2",
				Telephone:   "13631210111",
				Email:       "29400123@qq.com",
				Nickname:    "TEST03",
				Avatar:      "",
				Description: "Never Settle ..",
				Sex:         "0",
				Birthday:    1578903121862,
				Location:    "",
			}

			err := u.FindOneAndUpdateProfile(userCriteria, utils.TransformStructToMap(newProfile))
			user, _ := u.FindByUserId(newProfile.UserId)

			So(err, ShouldBeNil)
			So(user.Telephone, ShouldEqual, "13631210111")
			So(user.Email, ShouldEqual, "29400123@qq.com")
			So(user.Nickname, ShouldEqual, "TEST03")
		})
	})
	_ = u.RemoveUserByUserId(user1.UserId, true) // remove demo user
}
