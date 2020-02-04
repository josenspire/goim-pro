package authsrv

import (
	"context"
	. "github.com/smartystreets/goconvey/convey"
	protos "goim-pro/api/protos/salty"
	"goim-pro/pkg/utils"
	"testing"
)

func Test_smsServer_ObtainSMSCode(t *testing.T) {
	type args struct {
		ctx context.Context
		req *protos.GrpcResp
	}
	smsReq1 := &protos.ObtainSMSCodeReq{
		CodeType:  0,
		Telephone: "13631210000",
	}
	smsReq2 := &protos.ObtainSMSCodeReq{
		CodeType:  3,
		Telephone: "13631210000",
	}
	ctx := context.Background()
	any1, _ := utils.MarshalMessageToAny(smsReq1)
	req1 := &protos.GrpcReq{
		Data: any1,
	}
	any2, _ := utils.MarshalMessageToAny(smsReq2)
	req2 := &protos.GrpcReq{
		Data: any2,
	}
	Convey("obtain_sms_code", t, func() {
		Convey("testing_for_obtain_sms_code_of_register", func() {
			s := &smsServer{}
			gotRes, err := s.ObtainSMSCode(ctx, req1)
			if err != nil {
				t.Errorf("ObtainSMSCode() error = %v", err)
				return
			}
			So(gotRes.GetMessage(), ShouldEqual, "sending sms code success: 123456")
		})
		Convey("testing_for_invalid_code_type", func() {
			s := &smsServer{}
			gotRes, err := s.ObtainSMSCode(ctx, req2)
			if err != nil {
				t.Errorf("ObtainSMSCode() error = %v", err)
				return
			}
			So(gotRes.GetMessage(), ShouldEqual, "invalid request code type")
		})
	})
}
