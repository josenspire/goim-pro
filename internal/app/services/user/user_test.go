package usersrv

import (
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
	protos "goim-pro/api/protos/salty"
	mysqlsrv "goim-pro/pkg/db/mysql"
	"goim-pro/pkg/utils"
	"reflect"
	"testing"
)

var reqData1 = &protos.UserProfile{
	Telephone:   "13631210001",
	Email:       "123@qq.com",
	Username:    "13631210001",
	Nickname:    "JAMES001",
	Avatar:      "http://www.avatar.goo/123.png",
	Description: "Never Settle",
	Sex:         0,
	Birthday:    0,
	Location:    "CHINA-ZHA",
}

var reqData2 = &protos.UserProfile{
	Telephone:   "13631210002",
	Email:       "1233456@qq.com",
	Username:    "13631210002",
	Nickname:    "JAMES002",
	Avatar:      "http://www.avatar.goo/1234365.png",
	Description: "Never Settle 222",
	Sex:         1,
	Birthday:    0,
	Location:    "CHINA-ZHA-NR",
}

type MockUserRepo struct {
	mock.Mock
}

func initialMysqlDB() {
	_ = mysqlsrv.NewMysqlConnection().Connect()
}

func (m *MockUserRepo) IsTelephoneRegistered(telephone string) (bool, error) {
	args := m.Called(telephone)
	return args.Bool(0), args.Error(1)
}

func Test_Register(t *testing.T) {
	// TODO: should update the mock test process
	m := &MockUserRepo{}
	m.On("IsTelephoneRegistered", "13631210001").Return(true, nil)

	Convey("testing_grpc_user_register", t, func() {
		var ctx context.Context
		var req *protos.GrpcReq
		us := New()
		userReq := &protos.RegisterReq{
			RegisterType:     0,
			VerificationCode: "123456",
			UserProfile: &protos.UserProfile{
				Telephone:   "13631210001",
				Email:       "123@qq.com",
				Username:    "13631210001",
				Nickname:    "JAMES001",
				Avatar:      "http://www.avatar.goo/123.png",
				Description: "Never Settle",
				Sex:         0,
				Birthday:    utils.MakeTimestamp(),
				Location:    "CHINA-ZHA",
			},
			Password: "1234567890",
		}
		any, _ := utils.MarshalMessageToAny(userReq)
		req = &protos.GrpcReq{
			Data: any,
		}
		actualResp, err := us.Register(ctx, req)
		Convey("user_multiple_registration", func() {
			So(err, ShouldBeNil)
			So(actualResp.GetMessage(), ShouldEqual, "this telephone has been registered, please login")
		})
	})
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

func Test_userServer_Register(t *testing.T) {
	initialMysqlDB()

	Convey("testing_grpc_user_register", t, func() {
		var ctx context.Context
		var req *protos.GrpcReq
		us := New()
		userReq := &protos.RegisterReq{
			RegisterType:     0,
			VerificationCode: "123456",
			UserProfile: &protos.UserProfile{
				Telephone:   "13631210001",
				Email:       "123@qq.com",
				Username:    "13631210001",
				Nickname:    "JAMES001",
				Avatar:      "http://www.avatar.goo/123.png",
				Description: "Never Settle",
				Sex:         0,
				Birthday:    utils.MakeTimestamp(),
				Location:    "CHINA-ZHA",
			},
			Password: "1234567890",
		}
		any, _ := utils.MarshalMessageToAny(userReq)
		req = &protos.GrpcReq{
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
