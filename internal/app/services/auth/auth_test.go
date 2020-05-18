package authsrv

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	protos "goim-pro/api/protos/salty"
	"goim-pro/internal/app/repos/user"
	mysqlsrv "goim-pro/pkg/db/mysql"
	redsrv "goim-pro/pkg/db/redis"
	"testing"
)

func Test_authServer_ObtainSMSCode(t *testing.T) {
	_ = redsrv.NewRedisConnection().Connect()
	_ = mysqlsrv.NewMysqlConnection().Connect()

	m := &user.MockUserRepo{}
	m.On("IsTelephoneOrEmailRegistered", "13631210001", "").Return(false, nil)
	m.On("IsTelephoneOrEmailRegistered", "13631210002", "").Return(true, nil)
	m.On("IsTelephoneOrEmailRegistered", "13631210003", "").Return(true, nil)
	m.On("IsTelephoneOrEmailRegistered", "13631210004", "").Return(false, nil)

	authServer := New()
	userRepo = m

	Convey("testing_ObtainSMSCodeReq", t, func() {
		Convey("should_return_correct_sms_code_for_register", func() {
			actualCode, err := authServer.ObtainSMSCode("13631210001", protos.ObtainSMSCodeReq_REGISTER)

			fmt.Println(actualCode)
			str := myRedis.Get(fmt.Sprintf("%d-%s", protos.ObtainSMSCodeReq_REGISTER, "13631210001"))

			So(err, ShouldBeNil)
			So(actualCode, ShouldEqual, "123401")
			So(str, ShouldEqual, "123401")
		})
	})
}
