package authsrv

import (
	"context"
	"encoding/json"
	"github.com/golang/protobuf/ptypes/any"
	. "github.com/smartystreets/goconvey/convey"
	"goim-pro/api/protos"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want protos.SMSServiceServer
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_smsServer_ObtainSMSCode(t *testing.T) {
	type args struct {
		ctx context.Context
		req *protos.BasicClientRequest
	}
	dataBytes1, _ := json.Marshal(&protos.SMSReq{
		CodeType:  0,
		Telephone: "13631210000",
		Extension: nil,
	})
	dataBytes2, _ := json.Marshal(&protos.SMSReq{
		CodeType:  3,
		Telephone: "13631210000",
		Extension: nil,
	})
	ctx := context.Background()
	req1 := &protos.BasicClientRequest{
		Code: 0,
		Data: &any.Any{
			TypeUrl: "",
			Value:   dataBytes1,
		},
		Message: "",
	}
	req2 := &protos.BasicClientRequest{
		Code: 0,
		Data: &any.Any{
			TypeUrl: "",
			Value:   dataBytes2,
		},
		Message: "",
	}
	Convey("obtain_sms_code", t, func() {
		Convey("testing_for_obtain_sms_code_of_register", func() {
			s := &smsServer{}
			gotRes, err := s.ObtainSMSCode(ctx, req1)
			if err != nil {
				t.Errorf("ObtainSMSCode() error = %v", err)
				return
			}
			So(string(gotRes.GetData().Value), ShouldEqual, "123456")
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
