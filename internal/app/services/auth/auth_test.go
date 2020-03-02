package authsrv

import (
	"context"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	protos "goim-pro/api/protos/salty"
	usersrv "goim-pro/internal/app/services/user"
	"goim-pro/pkg/utils"
	"net/http"
	"testing"
)

func Test_smsServer_ObtainSMSCode(t *testing.T) {
	m := &usersrv.MockUserRepo{}
	m.On("IsTelephoneOrEmailRegistered", "13631210001", "").Return(false, nil)
	m.On("IsTelephoneOrEmailRegistered", "13631210002", "").Return(true, nil)
	m.On("IsTelephoneOrEmailRegistered", "13631210003", "").Return(true, nil)
	m.On("IsTelephoneOrEmailRegistered", "13631210004", "").Return(false, nil)

	Convey("testing_ObtainSMSCodeReq", t, func() {
		var ctx context.Context
		var req *protos.GrpcReq

		smsServer := &smsServer{
			userRepo: m,
		}
		Convey("should_return_correct_sms_code_for_register", func() {
			smsReq := &protos.ObtainSMSCodeReq{
				CodeType:  protos.ObtainSMSCodeReq_REGISTER,
				Telephone: "13631210001",
			}
			any, _ := utils.MarshalMessageToAny(smsReq)
			req = &protos.GrpcReq{
				Data: any,
			}
			actualResp, err := smsServer.ObtainSMSCode(ctx, req)

			fmt.Println(actualResp.GetMessage())

			str := myRedis.Get(fmt.Sprintf("%d-%s", protos.ObtainSMSCodeReq_REGISTER, "13631210001"))

			So(err, ShouldBeNil)
			So(actualResp.Code, ShouldEqual, http.StatusOK)
			So(len(str.String()), ShouldEqual, 6)
		})
	})
}
