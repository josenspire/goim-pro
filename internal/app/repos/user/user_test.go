package user

import (
	"github.com/jinzhu/gorm"
	"goim-pro/config"
	"goim-pro/internal/app/models"
	mysqlsrv "goim-pro/internal/db/mysql"
	"goim-pro/pkg/utils"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var user1 = &models.User{
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

var user2 = &models.User{
	Password: "1234567890",
	UserProfile: models.UserProfile{
		UserId:      "3",
		Telephone:   "13631210022",
		Email:       "294001@qq.com",
		Nickname:    "TEST02",
		Description: "Never Settle",
	},
}

var db *gorm.DB

func init() {
	mysqlsrv.NewMysql()
	db = mysqlsrv.GetMysql()
}

func TestUser_IsTelephoneRegistered(t *testing.T) {
	userRepo := NewUserRepo(db)
	_ = userRepo.Register(user2) // create a user
	Convey("Test_IsTelephoneRegistered", t, func() {
		Convey("Test_return_FALSE_with_exist_telephone", func() {
			isExist, err := userRepo.IsTelephoneOrEmailRegistered("13631210033", "123@qq.com")
			So(err, ShouldBeNil)
			So(isExist, ShouldBeFalse)
		})
		Convey("Test_return_FALSE_with_exist_email", func() {
			isExist, err := userRepo.IsTelephoneOrEmailRegistered("13631210044", "294001@qq.com")
			So(err, ShouldBeNil)
			So(isExist, ShouldBeTrue)
		})
		Convey("Test_return_TRUE", func() {
			isExist, err := userRepo.IsTelephoneOrEmailRegistered("13631210022", "")
			So(err, ShouldBeNil)
			So(isExist, ShouldBeTrue)
		})
	})
	_ = userRepo.RemoveUserByUserId(user2.UserId, true) // remove demo user
}

func TestUser_Register(t *testing.T) {
	userRepo := NewUserRepo(db)
	Convey("Test_Register", t, func() {
		Convey("Registration_successful", func() {
			err := userRepo.Register(user1)
			So(err, ShouldBeNil)
		})
		Convey("Registration_fail_by_telephone_exist", func() {
			err := userRepo.Register(user1)
			So(err, ShouldNotBeNil)
		})
	})

	_ = userRepo.RemoveUserByUserId(user1.UserId, true)
}

func TestUser_RemoveUserByUserId(t *testing.T) {
	userRepo := NewUserRepo(db)

	Convey("TestUserRepo_RemoveUserByUserID", t, func() {
		Convey("testing_RemoveUserByUserID_success", func() {
			err := userRepo.RemoveUserByUserId("1", false)
			So(err, ShouldBeNil)
		})
	})
}

func TestUser_Login(t *testing.T) {
	userRepo := NewUserRepo(db)

	_ = userRepo.Register(user2) // create a user

	Convey("TestUserRepo_LoginByTelephone", t, func() {
		Convey("login_fail_with_incorrect_telephone_and_password", func() {
			enPassword := utils.NewSHA256("12345678901", config.GetApiSecretKey())
			user, err := userRepo.QueryByTelephoneAndPassword("13631210022", enPassword)
			So(err, ShouldBeNil)
			So(user, ShouldBeNil)
		})
		Convey("login_success_then_return_userProfile", func() {
			enPassword := utils.NewSHA256("1234567890", config.GetApiSecretKey())
			currUser, err := userRepo.QueryByTelephoneAndPassword("13631210022", enPassword)
			So(err, ShouldBeNil)
			So(currUser.UserId, ShouldEqual, "3")
		})
	})
	_ = userRepo.RemoveUserByUserId(user2.UserId, true) // remove demo user
}

func TestUser_LoginByEmail(t *testing.T) {
	userRepo := NewUserRepo(db)

	_ = userRepo.Register(user2) // create a user
	Convey("TestUserRepo_LoginByEmail", t, func() {
		Convey("login_fail_with_incorrect_email_and_password", func() {
			enPassword := utils.NewSHA256("12345678901", config.GetApiSecretKey())
			user, err := userRepo.QueryByEmailAndPassword("294001@qq.com", enPassword)
			So(err, ShouldBeNil)
			So(user, ShouldBeNil)
		})
		Convey("login_success_then_return_userProfile", func() {
			enPassword := utils.NewSHA256("1234567890", config.GetApiSecretKey())
			currUser, err := userRepo.QueryByEmailAndPassword("294001@qq.com", enPassword)
			So(err, ShouldBeNil)
			So(currUser.UserId, ShouldEqual, "3")
			So(currUser.Email, ShouldEqual, "294001@qq.com")
		})
	})
	_ = userRepo.RemoveUserByUserId(user2.UserId, true) // remove demo user
}

func TestUser_ResetPasswordByTelephone(t *testing.T) {
	userRepo := NewUserRepo(db)

	_ = userRepo.Register(user1) // create a user
	Convey("TestUserRepo_ResetPasswordByTelephone", t, func() {
		Convey("update_successful_by_telephone", func() {
			var telephone string = "13631210010"
			var password string = "111111111"
			newPassword := utils.NewSHA256(password, config.GetApiSecretKey())

			err := userRepo.ResetPasswordByTelephone(telephone, newPassword)
			So(err, ShouldBeNil)

			user, err := userRepo.QueryByTelephoneAndPassword(telephone, newPassword)
			So(err, ShouldBeNil)
			So(user.Telephone, ShouldEqual, telephone)
		})
	})
	_ = userRepo.RemoveUserByUserId(user1.UserId, true) // remove demo user
}

func TestUser_ResetPasswordByEmail(t *testing.T) {
	userRepo := NewUserRepo(db)

	_ = userRepo.Register(user1) // create a user
	Convey("TestUserRepo_ResetPasswordByEmail", t, func() {
		Convey("update_successful_by_email", func() {
			var email string = "294001@qq.com"
			var password string = "2222222222"
			newPassword := utils.NewSHA256(password, config.GetApiSecretKey())

			err := userRepo.ResetPasswordByEmail(email, newPassword)
			So(err, ShouldBeNil)

			user, err := userRepo.QueryByEmailAndPassword(email, newPassword)
			So(err, ShouldBeNil)
			So(user.Email, ShouldEqual, email)
		})
	})
	_ = userRepo.RemoveUserByUserId(user1.UserId, true) // remove demo user
}

func TestUser_GetUserByUserId(t *testing.T) {
	userRepo := NewUserRepo(db)

	_ = userRepo.Register(user1) // create a user
	Convey("TestUserRepo_GetUserByUserId", t, func() {
		Convey("get_user_info_success_then_return", func() {
			var userId = "2"

			user, err := userRepo.FindByUserId(userId)
			So(err, ShouldBeNil)
			So(user.Email, ShouldEqual, "294001@qq.com")
		})
		Convey("get_user_info_failed_by_invalid_userId", func() {
			var userId = "3"

			user, err := userRepo.FindByUserId(userId)
			So(err, ShouldBeNil)
			So(user, ShouldBeNil)
		})
	})
	_ = userRepo.RemoveUserByUserId(user1.UserId, true) // remove demo user
}

func TestUser_FindOneUser(t *testing.T) {
	userRepo := NewUserRepo(db)

	_ = userRepo.Register(user1) // create a user
	Convey("TestUserRepo_FindOneUser", t, func() {
		Convey("get_user_info_success_by_telephone_then_return", func() {
			var userCriteria = map[string]interface{}{
				"Telephone": "13631210010",
			}

			user, err := userRepo.FindOneUser(userCriteria)
			So(err, ShouldBeNil)
			So(user.Email, ShouldEqual, "294001@qq.com")
		})
		Convey("get_user_info_success_by_telephone_and_email_then_return", func() {
			var userCriteria = map[string]interface{}{
				"Email":     "294001@qq.com",
				"Telephone": "13631210010",
			}

			user, err := userRepo.FindOneUser(userCriteria)
			So(err, ShouldBeNil)
			So(user.Email, ShouldEqual, "294001@qq.com")
		})
		Convey("get_user_info_failed_by_invalid_email_then_return_nil", func() {
			var userCriteria = map[string]interface{}{
				"Email":     "294001@qq.com",
				"Telephone": "13631210015",
			}

			user, err := userRepo.FindOneUser(userCriteria)
			So(err, ShouldBeNil)
			So(user, ShouldBeNil)
		})
	})
	_ = userRepo.RemoveUserByUserId(user1.UserId, true) // remove demo user
}

func TestUser_FindOneAndUpdateProfile(t *testing.T) {
	userRepo := NewUserRepo(db)

	_ = userRepo.Register(user1) // create a user
	Convey("TestUserRepo_FindOneAndUpdateProfile", t, func() {
		Convey("find_user_by_userId_and_update_profile_success", func() {
			var userCriteria = map[string]interface{}{
				"UserId": "2",
			}
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
			err := userRepo.FindOneAndUpdateProfile(userCriteria, utils.TransformStructToMap(newProfile))
			user, _ := userRepo.FindByUserId(newProfile.UserId)

			So(err, ShouldBeNil)
			So(user.Telephone, ShouldEqual, "13631210111")
			So(user.Email, ShouldEqual, "29400123@qq.com")
			So(user.Nickname, ShouldEqual, "TEST03")
		})
	})
	_ = userRepo.RemoveUserByUserId(user1.UserId, true) // remove demo user
}
