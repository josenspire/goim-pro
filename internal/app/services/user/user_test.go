package usersrv

import (
	"context"
	. "github.com/smartystreets/goconvey/convey"
	protos "goim-pro/api/protos/salty"
	"goim-pro/internal/app/repos/user"
	"goim-pro/pkg/utils"
	"testing"
)

var pbUserProfile1 = &protos.UserProfile{
	UserID:      "13631210001",
	Telephone:   "13631210001",
	Email:       "123@qq.com",
	Username:    "13631210001",
	Nickname:    "JAMES001",
	Avatar:      "http://www.avatar.goo/123.png",
	Description: "Never Settle",
	Sex:         0,
	Birthday:    utils.MakeTimestamp(),
	Location:    "CHINA-ZHA",
}

var pbUserProfile2 = &protos.UserProfile{
	UserID:      "13631210002",
	Telephone:   "13631210002",
	Email:       "12345@qq.com",
	Username:    "13631210001",
	Nickname:    "JAMES001",
	Avatar:      "http://www.avatar.goo/123.png",
	Description: "Never Settle",
	Sex:         0,
	Birthday:    utils.MakeTimestamp(),
	Location:    "CHINA-ZHA",
}

var modelUserProfile1 = &user.UserProfile{
	UserID:      "13631210001",
	Telephone:   "13631210001",
	Email:       "123@qq.com",
	Username:    "13631210001",
	Nickname:    "JAMES001",
	Avatar:      "http://www.avatar.goo/123.png",
	Description: "Never Settle",
	Sex:         "MALE",
	Birthday:    utils.MakeTimestamp(),
	Location:    "CHINA-ZHA",
}
var modelUserProfile2 = &user.UserProfile{
	UserID:      "13631210002",
	Telephone:   "13631210002",
	Email:       "12345@qq.com",
	Username:    "13631210001",
	Nickname:    "JAMES001",
	Avatar:      "http://www.avatar.goo/123.png",
	Description: "Never Settle",
	Sex:         "MALE",
	Birthday:    utils.MakeTimestamp(),
	Location:    "CHINA-ZHA",
}

func Test_Register(t *testing.T) {
	m := &MockUserRepo{}
	m.On("IsTelephoneOrEmailRegistered", "13631210001", "123@qq.com").Return(true, nil)
	m.On("IsTelephoneOrEmailRegistered", "13631210002", "12345@qq.com").Return(false, nil)
	m.On("Register", &user.User{
		Password:    "1234567890",
		UserProfile: *modelUserProfile1,
	}).Return(nil)
	m.On("Register", &user.User{
		Password:    "1234567890",
		UserProfile: *modelUserProfile2,
	}).Return(nil)

	Convey("testing_grpc_user_register", t, func() {
		var ctx context.Context
		var req *protos.GrpcReq
		us := &userService{
			userRepo: m,
		}
		Convey("user_registration_fail_by_exist_telephone", func() {
			userReq := &protos.RegisterReq{
				RegisterType:     0,
				VerificationCode: "123456",
				UserProfile:      pbUserProfile1,
				Password:         "1234567890",
			}
			any, _ := utils.MarshalMessageToAny(userReq)
			req = &protos.GrpcReq{
				Data: any,
			}
			actualResp, err := us.Register(ctx, req)
			So(err, ShouldBeNil)
			So(actualResp.GetMessage(), ShouldEqual, "the telephone or email has been registered, please login")
		})
		Convey("user_registration_successful", func() {
			userReq := &protos.RegisterReq{
				RegisterType:     0,
				VerificationCode: "123456",
				UserProfile:      pbUserProfile2,
				Password:         "1234567890",
			}
			any, _ := utils.MarshalMessageToAny(userReq)
			req = &protos.GrpcReq{
				Data: any,
			}
			actualResp, err := us.Register(ctx, req)

			registerResp := &protos.RegisterResp{}
			err = utils.UnMarshalAnyToMessage(actualResp.GetData(), registerResp)

			So(err, ShouldBeNil)
			So(actualResp.GetMessage(), ShouldEqual, "user registration successful")
			So(registerResp.Profile.GetTelephone(), ShouldEqual, pbUserProfile2.GetTelephone())
		})
	})
}
