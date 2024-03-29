package usersrv

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
	protos "goim-pro/api/protos/salty"
	"goim-pro/config"
	consts "goim-pro/internal/app/constants"
	"goim-pro/internal/app/models"
	"goim-pro/internal/app/repos/user"
	mysqlsrv "goim-pro/internal/db/mysql"
	redsrv "goim-pro/internal/db/redis"
	"goim-pro/pkg/errors"
	"goim-pro/pkg/utils"
	"testing"
)

var pbUserProfile1 = &protos.UserProfile{
	UserId:      "13631210001",
	Telephone:   "13631210001",
	Email:       "123@qq.com",
	Nickname:    "JAMES001",
	Avatar:      "http://www.avatar.goo/123.png",
	Description: "Never Settle",
	Sex:         0,
	Birthday:    utils.MakeTimestamp(),
	Location:    "CHINA-ZHA",
}

var pbUserProfile2 = &protos.UserProfile{
	UserId:      "13631210002",
	Telephone:   "13631210002",
	Email:       "12345@qq.com",
	Nickname:    "JAMES001",
	Avatar:      "http://www.avatar.goo/123.png",
	Description: "Never Settle",
	Sex:         0,
	Birthday:    utils.MakeTimestamp(),
	Location:    "CHINA-ZHA",
}

var modelUserProfile1 = &models.UserProfile{
	UserId:      "13631210001",
	Telephone:   "13631210001",
	Email:       "123@qq.com",
	Nickname:    "JAMES001",
	Avatar:      "http://www.avatar.goo/123.png",
	Description: "Never Settle",
	Sex:         "MALE",
	Birthday:    utils.MakeTimestamp(),
	Location:    "CHINA-ZHA",
}
var modelUserProfile2 = &models.UserProfile{
	UserId:      "13631210002",
	Telephone:   "13631210002",
	Email:       "12345@qq.com",
	Nickname:    "JAMES001",
	Avatar:      "http://www.avatar.goo/123.png",
	Description: "Never Settle",
	Sex:         "MALE",
	Birthday:    utils.MakeTimestamp(),
	Location:    "CHINA-ZHA",
}

var modelUser1 = &models.User{
	DeviceId:    "XIAOMI10",
	OsVersion:   "1",
	UserProfile: *modelUserProfile1,
}
var modelUser2 = &models.User{
	DeviceId:    "OnePlus 8",
	OsVersion:   "1",
	UserProfile: *modelUserProfile2,
}

func init() {
	mysqlsrv.NewMysql()
}

func Test_Register(t *testing.T) {
	m := &user.MockUserRepo{}
	m.On("IsTelephoneOrEmailRegistered", "13631210001", "123@qq.com").Return(true, nil)
	m.On("IsTelephoneOrEmailRegistered", "13631210002", "12345@qq.com").Return(false, nil)
	m.On("Register", mock.Anything).Return(nil)

	var registerKey1 = fmt.Sprintf("%d-%s", consts.CodeTypeRegister, "13631210001")
	var registerKey2 = fmt.Sprintf("%d-%s", consts.CodeTypeRegister, "13631210002")
	r := new(redsrv.MockCmdable)
	r.On("RGet", registerKey1).Return("123456")
	r.On("RGet", registerKey2).Return("123456")
	r.On("RDel", registerKey1).Return(0)
	r.On("RDel", registerKey2).Return(0)

	us := New()
	userRepo = m
	myRedis = r

	Convey("testing_grpc_user_register", t, func() {
		Convey("user_registration_fail_by_exist_telephone", func() {
			tErr := us.Register(pbUserProfile1, "1234567890", "LOCAL_DEV", 1)
			So(tErr, ShouldNotBeNil)
			So(tErr.Code, ShouldEqual, protos.StatusCode_STATUS_ACCOUNT_EXISTS)
			So(tErr.Detail, ShouldEqual, errmsg.ErrAccountAlreadyExists.Error())
		})
		Convey("user_registration_successful", func() {
			tErr := us.Register(pbUserProfile2, "1234567890", "LOCAL_DEV", 1)
			So(tErr, ShouldBeNil)
		})
	})
}

func Test_userService_Login(t *testing.T) {

	m := &user.MockUserRepo{}
	enPassword := utils.NewSHA256("1234567890", config.GetApiSecretKey())
	m.On("QueryByTelephoneAndPassword", "13631210001", enPassword).Return(modelUser1, nil)
	m.On("QueryByEmailAndPassword", "123@qq.com", enPassword).Return(modelUser1, nil)
	m.On("FindOneUser", map[string]interface{}{"telephone": "13631210001"}).Return(modelUser1, nil)

	var tokenKey = "TK-13631210001"
	var loginKey1 = fmt.Sprintf("%d-%s", consts.CodeTypeLogin, "13631210001")
	var loginKey2 = fmt.Sprintf("%d-%s", consts.CodeTypeLogin, "13631210002")

	r := new(redsrv.MockCmdable)
	r.On("RGet", loginKey1).Return("123456")
	r.On("RGet", loginKey2).Return("123456")

	r.On("RSet", tokenKey, mock.Anything, consts.ThreeDays).Return(nil)
	//r.On("RSet", loginKey1, mock.Anything, consts.ThreeDays).Return("123456")

	r.On("RDel", loginKey1).Return(0)
	r.On("RDel", loginKey1).Return(0)

	us := New()
	userRepo = m
	myRedis = r

	deviceId := "XIAOMI10"
	osVersion := protos.GrpcReq_ANDROID

	Convey("Test_Login", t, func() {
		Convey("should_login_successful_by_telephone_and_password_then_return_user_profile", func() {
			user, token, tErr := us.Login("13631210001", "", enPassword, deviceId, osVersion)
			So(tErr, ShouldBeNil)
			So(user.Telephone, ShouldEqual, "13631210001")
			So(token, ShouldNotBeEmpty)
		})
		Convey("should_login_successful_by_email_and_password_then_return_user_profile", func() {
			user, token, tErr := us.Login("", "123@qq.com", enPassword, deviceId, osVersion)
			So(tErr, ShouldBeNil)
			So(user.Email, ShouldEqual, "123@qq.com")
			So(token, ShouldNotBeEmpty)
		})
		Convey("should_login_by_account_fail_when_deviceId_change_then_return_security_verification_require_message", func() {
			user, token, tErr := us.Login("13631210001", "", enPassword, "OnePlus 8", osVersion)
			So(token, ShouldBeEmpty)
			So(user, ShouldBeNil)
			So(tErr, ShouldNotBeNil)
			So(tErr.Code, ShouldEqual, protos.StatusCode_STATUS_ACCOUNT_AUTHORIZED_REQUIRED)
		})
	})
}

func Test_userService_ResetPassword(t *testing.T) {

	//telephone, email := "13631210001", "123@qq.com"
	enPassword := utils.NewSHA256("1234567890", config.GetApiSecretKey())
	errEnPassword := utils.NewSHA256("123456789111", config.GetApiSecretKey())
	newEnPassword := utils.NewSHA256("1122334455", config.GetApiSecretKey())

	m := &user.MockUserRepo{}
	m.On("IsTelephoneOrEmailRegistered", "13631210001", "").Return(true, nil)
	m.On("IsTelephoneOrEmailRegistered", "", "123@qq.com").Return(true, nil)
	m.On("IsTelephoneOrEmailRegistered", "", "123456@qq.com").Return(false, errmsg.ErrUserNotExists)

	m.On("QueryByTelephoneAndPassword", "13631210001", enPassword).Return(modelUser1, nil)
	m.On("QueryByEmailAndPassword", "123@qq.com", enPassword).Return(modelUser1, nil)
	m.On("QueryByTelephoneAndPassword", "13631210001", errEnPassword).Return(nil, errmsg.ErrAccountOrPwdInvalid)

	m.On("ResetPasswordByTelephone", "13631210001", newEnPassword).Return(nil)
	m.On("ResetPasswordByEmail", "123@qq.com", newEnPassword).Return(nil)

	var registerKey1 = fmt.Sprintf("%d-%s", consts.CodeTypeResetPassword, "13631210001")
	r := new(redsrv.MockCmdable)
	r.On("RGet", registerKey1).Return("123456")
	r.On("RDel", registerKey1).Return(0)

	us := New()
	userRepo = m
	myRedis = r

	Convey("Test_ResetPassword", t, func() {
		Convey("user_reset_password_successful_by_telephone_with_old_password", func() {
			tErr := us.ResetPassword("13631210001", "", "1122334455", "LOCAL_DEV", 1)
			So(tErr, ShouldBeNil)
		})
		//Convey("failed_by_newPassword_same_as_old", func() {
		//	tErr := us.ResetPassword("13631210001", "", "1122334455", "LOCAL_DEV", 1)
		//	So(tErr, ShouldNotBeNil)
		//	So(tErr.Code, ShouldEqual, protos.StatusCode_STATUS_BAD_REQUEST)
		//	So(tErr.Detail, ShouldEqual, errmsg.ErrRepeatPassword.Error())
		//})
		Convey("failed_by_not_exist_account", func() {
			tErr := us.ResetPassword("", "123456@qq.com", "1122334455", "LOCAL_DEV", 1)
			So(tErr, ShouldNotBeNil)
			So(tErr.Code, ShouldEqual, protos.StatusCode_STATUS_BAD_REQUEST)
			So(tErr.Detail, ShouldEqual, errmsg.ErrUserNotExists.Error())
		})
		//Convey("failed_by_invalid_oldPassword", func() {
		//	tErr := us.ResetPassword("13631210001", "", "1122334455", "LOCAL_DEV", 1)
		//	So(tErr, ShouldNotBeNil)
		//	So(tErr.Detail, ShouldEqual, errmsg.ErrAccountOrPwdInvalid.Error())
		//})
	})
}

func Test_userService_GetUserInfo(t *testing.T) {

	m := &user.MockUserRepo{}
	m.On("FindByUserId", "13631210001").Return(modelUser1, nil)
	m.On("FindByUserId", "13631210002").Return(nil, nil)

	us := New()
	userRepo = m

	Convey("Test_GetUserInfo", t, func() {
		Convey("should_return_user_when_given_correct_userId", func() {
			profile, tErr := us.GetUserInfo("13631210001")
			So(tErr, ShouldBeNil)
			So(profile.Telephone, ShouldEqual, "13631210001")
		})
		Convey("should_get_user_info_fail_by_incorrect_user_then_return_error", func() {
			profile, tErr := us.GetUserInfo("13631210002")
			So(profile, ShouldBeNil)
			So(tErr, ShouldNotBeNil)
			So(tErr.Code, ShouldEqual, protos.StatusCode_STATUS_BAD_REQUEST)
			So(tErr.Detail, ShouldEqual, errmsg.ErrInvalidUserId.Error())
		})
	})
}

func Test_userService_QueryUserInfo(t *testing.T) {
	userCriteria1 := map[string]interface{}{
		"telephone": "13631210001",
	}

	userCriteria2 := map[string]interface{}{
		"email": "123@qq.com",
	}

	userCriteria3 := map[string]interface{}{
		"telephone": "13631210012",
	}

	m := &user.MockUserRepo{}
	m.On("FindOneUser", userCriteria1).Return(modelUser1, nil)
	m.On("FindOneUser", userCriteria2).Return(modelUser1, nil)
	m.On("FindOneUser", userCriteria3).Return(nil, nil)

	us := New()
	userRepo = m

	Convey("Testing_QueryUserInfo", t, func() {
		Convey("should_return_user_when_given_correct_telephone", func() {
			profile, tErr := us.QueryUserInfo("13631210001", "")
			So(tErr, ShouldBeNil)
			So(profile.Telephone, ShouldEqual, "13631210001")
		})
		Convey("should_return_user_when_given_correct_email", func() {
			profile, tErr := us.QueryUserInfo("", "123@qq.com")
			So(tErr, ShouldBeNil)
			So(profile.Telephone, ShouldEqual, "13631210001")
			So(profile.Email, ShouldEqual, "123@qq.com")
		})
		Convey("should_return_err_when_given_un_exists_telephone", func() {
			profile, tErr := us.QueryUserInfo("13631210012", "")

			So(profile, ShouldBeNil)
			So(tErr, ShouldNotBeNil)
			So(tErr.Code, ShouldEqual, protos.StatusCode_STATUS_BAD_REQUEST)
			So(tErr.Detail, ShouldEqual, errmsg.ErrUserNotExists.Error())
		})
	})
}

func Test_userService_UpdateUserInfo(t *testing.T) {

	criteria1 := &models.User{}
	criteria1.UserId = "13631210001"

	criteria2 := &models.User{}
	criteria2.UserId = "13631210002"

	newProfile1 := models.UserProfile{
		UserId:      "13631210001",
		Telephone:   "13631214444",
		Email:       "123456@qq.com",
		Nickname:    "JAMESYANG01",
		Avatar:      "",
		Description: "",
		Sex:         "0",
		Birthday:    0,
		Location:    "ZHA",
	}
	newProfile2 := models.UserProfile{
		UserId:      "13631210002",
		Telephone:   "13631214444",
		Email:       "123456@qq.com",
		Nickname:    "JAMESYANG01",
		Avatar:      "",
		Description: "",
		Sex:         "0",
		Birthday:    0,
		Location:    "ZHA",
	}

	m := &user.MockUserRepo{}
	m.On("FindByUserId", "13631210001").Return(modelUser1, nil)
	m.On("FindByUserId", "13631210003").Return(nil, nil)
	m.On("FindOneAndUpdateProfile", criteria1, utils.TransformStructToMap(newProfile1)).Return(nil)
	m.On("FindOneAndUpdateProfile", criteria2, utils.TransformStructToMap(newProfile2)).Return(nil)

	us := New()
	userRepo = m

	Convey("Testing_UpdateUserInfo", t, func() {
		Convey("should_update_profile_success_when_given_correct_newProfile", func() {
			tErr := us.UpdateUserInfo("13631210001", modelUserProfile1)
			So(tErr, ShouldBeNil)
		})

		Convey("should_update_profile_failed_when_given_incorrect_profile_with_userId", func() {
			tErr := us.UpdateUserInfo("13631210003", modelUserProfile2)
			So(tErr, ShouldNotBeNil)
			So(tErr.Detail, ShouldEqual, errmsg.ErrInvalidUserId.Error())
		})
	})
}
