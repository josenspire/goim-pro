package usersrv

import (
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"goim-pro/api/protos"
	"goim-pro/pkg/utils"
	"reflect"
	"testing"
)

var reqData1 = &protos.UserProfile{
	Telephone: "13631210001",
	Email:     "123@qq.com",
	Username:  "13631210001",
	Nickname:  "JAMES001",
	Avatar:    "http://www.avatar.goo/123.png",
	Signature: "Never Settle",
	Sex:       0,
	Birthday:  nil,
	Location:  "CHINA-ZHA",
}

var reqData2 = &protos.UserProfile{
	Telephone: "13631210002",
	Email:     "1233456@qq.com",
	Username:  "13631210002",
	Nickname:  "JAMES002",
	Avatar:    "http://www.avatar.goo/1234365.png",
	Signature: "Never Settle 222",
	Sex:       1,
	Birthday:  nil,
	Location:  "CHINA-ZHA-NR",
}

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want protos.UserServiceServer
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

func Test_userServer_Login(t *testing.T) {
	type args struct {
		ctx context.Context
		req *protos.BasicClientRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *protos.BasicServerResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := &userServer{}
			got, err := us.Login(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Login() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userServer_Register(t *testing.T) {
	Convey("testing_grpc_user_register", t, func() {
		var ctx context.Context
		var req *protos.BasicClientRequest
		us := &userServer{}
		userReq := protos.UserReq{
			CodeType:         0,
			VerificationCode: "123456",
			UserProfile: &protos.UserProfile{
				Telephone: "13631210001",
				Email:     "123@qq.com",
				Username:  "13631210001",
				Nickname:  "JAMES001",
				Avatar:    "http://www.avatar.goo/123.png",
				Signature: "Never Settle",
				Sex:       0,
				Birthday:  nil,
				Location:  "CHINA-ZHA",
			},
			Password: "1234567890",
		}
		any, _ := utils.NewMarshalAny(&userReq)
		req = &protos.BasicClientRequest{
			Data: any,
		}
		actualResp, err := us.Register(ctx, req)
		Convey("user_registration_successful", func() {
			So(err, ShouldBeNil)
			So(actualResp.GetMessage(), ShouldEqual, "user registration successful")
		})
		Convey("user_multiple_registration", func() {
			So(err, ShouldBeNil)
			So(actualResp.GetMessage(), ShouldEqual, "this telephone has been registered, please login")
		})
	})
}
