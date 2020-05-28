package authsrv

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	protos "goim-pro/api/protos/salty"
	"goim-pro/internal/app/constants"
	"goim-pro/internal/app/repos/user"
	mysqlsrv "goim-pro/pkg/db/mysql"
	redsrv "goim-pro/pkg/db/redis"
	"testing"
)

func Test_authServer_ObtainSMSCode(t *testing.T) {
	_ = mysqlsrv.NewMysql()

	m := &user.MockUserRepo{}
	m.On("IsTelephoneOrEmailRegistered", "13631210001", "").Return(false, nil)
	m.On("IsTelephoneOrEmailRegistered", "13631210002", "").Return(true, nil)
	m.On("IsTelephoneOrEmailRegistered", "13631210003", "").Return(true, nil)
	m.On("IsTelephoneOrEmailRegistered", "13631210004", "").Return(false, nil)

	var registerKey = fmt.Sprintf("%d-%s", protos.ObtainSMSCodeReq_REGISTER, "13631210001")

	r := new(redsrv.MockCmdable)
	r.On("RSet", registerKey, "123401", consts.MinuteOf15).Return(nil)
	r.On("RGet", registerKey).Return("123401")

	authServer := new(AuthService)
	userRepo = m
	myRedis = r

	Convey("testing_ObtainSMSCodeReq", t, func() {
		Convey("should_return_correct_sms_code_for_register", func() {
			actualCode, err := authServer.ObtainSMSCode("13631210001", protos.ObtainSMSCodeReq_REGISTER)

			fmt.Println(actualCode)

			So(err, ShouldBeNil)
			So(actualCode, ShouldEqual, "123401")
		})
	})
}
