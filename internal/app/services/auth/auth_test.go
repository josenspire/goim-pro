package authsrv

import (
	"fmt"
	protos "goim-pro/api/protos/salty"
	consts "goim-pro/internal/app/constants"
	"goim-pro/internal/app/repos/user"
	mysqlsrv "goim-pro/internal/db/mysql"
	redsrv "goim-pro/internal/db/redis"
	errmsg "goim-pro/pkg/errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	mysqlsrv.NewMysql()
}

func TestAuthService_ObtainSMSCode(t *testing.T) {
	m := &user.MockUserRepo{}
	m.On("IsTelephoneOrEmailRegistered", "13631210001", "").Return(false, nil)
	m.On("IsTelephoneOrEmailRegistered", "13631210002", "").Return(true, nil)
	m.On("IsTelephoneOrEmailRegistered", "13631210003", "").Return(false, nil)

	var registerKey = fmt.Sprintf("%d-%s", protos.SMSOperationType_REGISTER, "13631210001")
	var loginKey = fmt.Sprintf("%d-%s", protos.SMSOperationType_LOGIN, "13631210002")
	var resetPasswordKey = fmt.Sprintf("%d-%s", protos.SMSOperationType_RESET_PASSWORD, "13631210003")

	r := new(redsrv.MockCmdable)
	r.On("RSet", registerKey, "123401", consts.MinuteOf15).Return(nil)
	r.On("RSet", loginKey, "123402", consts.MinuteOf15).Return(nil)
	r.On("RSet", resetPasswordKey, "123403", consts.MinuteOf15).Return(nil)

	authServer := new(AuthService)
	userRepo = m
	myRedis = r

	Convey("testing_ObtainSMSCode", t, func() {
		Convey("should_return_correct_sms_code_for_register", func() {
			actualCode, err := authServer.ObtainSMSCode("13631210001", protos.SMSOperationType_REGISTER)

			fmt.Println(actualCode)

			So(err, ShouldBeNil)
			So(actualCode, ShouldEqual, "123401")
		})
		Convey("should_return_correct_sms_code_for_login", func() {
			actualCode, err := authServer.ObtainSMSCode("13631210002", protos.SMSOperationType_LOGIN)

			fmt.Println(actualCode)

			So(err, ShouldBeNil)
			So(actualCode, ShouldEqual, "123402")
		})
		Convey("should_not_return_reset_password_code_by_telephone_unExist_error", func() {
			actualCode, err := authServer.ObtainSMSCode("13631210003", protos.SMSOperationType_RESET_PASSWORD)

			fmt.Println(actualCode)

			So(err, ShouldNotBeNil)
			So(err.Detail, ShouldEqual, errmsg.ErrAccountNotExists.Error())
		})
	})
}

func TestAuthService_VerifySMSCode(t *testing.T) {
	m := &user.MockUserRepo{}
	m.On("IsTelephoneOrEmailRegistered", "13631210001", "").Return(true, nil)

	var registerKey = fmt.Sprintf("%d-%s", protos.SMSOperationType_REGISTER, "13631210001")

	r := new(redsrv.MockCmdable)
	r.On("RGet", registerKey).Return("123401")

	r.On("RDel", registerKey).Return(0)

	authServer := new(AuthService)
	userRepo = m
	myRedis = r

	Convey("testing_VerifySMSCode", t, func() {
		Convey("should_verify_password_when_given_correct_type_and_code_then_delete_the_code_cache", func() {
			isPass, err := authServer.VerifySMSCode("13631210001", protos.SMSOperationType_REGISTER, "123401", "LOCAL_DEV")

			So(err, ShouldBeNil)
			So(isPass, ShouldBeTrue)
		})
		Convey("should_verify_failed_when_given_incorrect_code", func() {
			isPass, err := authServer.VerifySMSCode("13631210001", protos.SMSOperationType_REGISTER, "123402", "LOCAL_DEV")

			So(isPass, ShouldBeFalse)
			So(err, ShouldNotBeNil)
			So(err.Detail, ShouldEqual, errmsg.ErrInvalidVerificationCode.Error())
		})
	})
}
