package grpc

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_isOnWhiteList(t *testing.T) {
	Convey("Testing_white_list_checking", t, func() {
		Convey("should_on_white_list", func() {
			var smsMethod = "/com.salty.protos.SMSService/ObtainSMSCode"
			isValid := isOnWhiteList(smsMethod)
			So(isValid, ShouldBeTrue)
		})
		Convey("should_not_on_white_list", func() {
			var registerMethod = "/com.salty.protos.UserService/Register"
			isValid := isOnWhiteList(registerMethod)
			So(isValid, ShouldBeFalse)
		})
	})
}
