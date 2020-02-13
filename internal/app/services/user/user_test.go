package usersrv

import (
	"context"
	protos "goim-pro/api/protos/salty"
	"goim-pro/config"
	"goim-pro/internal/app/repos/user"
	"goim-pro/pkg/utils"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var pbUserProfile1 = &protos.UserProfile{
	UserId:      "13631210001",
	Telephone:   "13631210001",
	Email:       "123@qq.com",
	Nickname:    "JAMES001",
	Avatar:      "http://www.avatar.goo/123.png",
	Description: "Never Settle",
	Sex:         0,
	Birthday:    utils.MakeTimestamp(),
	Location:    "CHINA-ZHA",
}

var pbUserProfile2 = &protos.UserProfile{
	UserId:      "13631210002",
	Telephone:   "13631210002",
	Email:       "12345@qq.com",
	Nickname:    "JAMES001",
	Avatar:      "http://www.avatar.goo/123.png",
	Description: "Never Settle",
	Sex:         0,
	Birthday:    utils.MakeTimestamp(),
	Location:    "CHINA-ZHA",
}

var modelUserProfile1 = &user.UserProfile{
	UserId:      "13631210001",
	Telephone:   "13631210001",
	Email:       "123@qq.com",
	Nickname:    "JAMES001",
	Avatar:      "http://www.avatar.goo/123.png",
	Description: "Never Settle",
	Sex:         "MALE",
	Birthday:    utils.MakeTimestamp(),
	Location:    "CHINA-ZHA",
}
var modelUserProfile2 = &user.UserProfile{
	UserId:      "13631210002",
	Telephone:   "13631210002",
	Email:       "12345@qq.com",
	Nickname:    "JAMES001",
	Avatar:      "http://www.avatar.goo/123.png",
	Description: "Never Settle",
	Sex:         "MALE",
	Birthday:    utils.MakeTimestamp(),
	Location:    "CHINA-ZHA",
}

var modelUser1 = &user.User{
	UserProfile: *modelUserProfile1,
}
var modelUser2 = &user.User{
	UserProfile: *modelUserProfile2,
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
				VerificationCode: "123456",
				Profile:          pbUserProfile1,
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
				VerificationCode: "123456",
				Profile:          pbUserProfile2,
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

func Test_userService_Login(t *testing.T) {
	// TODO：
}

func Test_userService_ResetPassword(t *testing.T) {
	//telephone, email := "13631210001", "123@qq.com"
	enPassword, _ := crypto.AESEncrypt("1234567890", config.GetApiSecretKey())
	newEnPassword, _ := crypto.AESEncrypt("1122334455", config.GetApiSecretKey())

	m := &MockUserRepo{}
	m.On("QueryByTelephoneAndPassword", "13631210001", enPassword).Return(modelUser1, nil)
	m.On("QueryByEmailAndPassword", "123@qq.com", enPassword).Return(modelUser1, nil)

	m.On("ResetPasswordByTelephone", "13631210001", newEnPassword).Return(nil)
	m.On("ResetPasswordByEmail", "123@qq.com", newEnPassword).Return(nil)

	Convey("testing_grpc_user_reset_password", t, func() {
		var ctx context.Context
		var req *protos.GrpcReq
		us := &userService{
			userRepo: m,
		}
		Convey("user_reset_password_successful_by_telephone", func() {
			resetPwdReq := &protos.ResetPasswordReq{
				OldPassword:      "1234567890",
				NewPassword:      "1122334455",
				VerificationCode: "112233",
				TargetAccount: &protos.ResetPasswordReq_Telephone{
					Telephone: "13631210001",
				},
			}
			any, _ := utils.MarshalMessageToAny(resetPwdReq)
			req = &protos.GrpcReq{
				Data: any,
			}
			actualResp, err := us.ResetPassword(ctx, req)
			So(err, ShouldBeNil)
			So(actualResp.GetCode(), ShouldEqual, 200)
		})
		Convey("user_reset_password_successful_by_email", func() {
			resetPwdReq := &protos.ResetPasswordReq{
				OldPassword:      "1234567890",
				NewPassword:      "1122334455",
				VerificationCode: "112233",
				TargetAccount: &protos.ResetPasswordReq_Email{
					Email: "123@qq.com",
				},
			}
			any, _ := utils.MarshalMessageToAny(resetPwdReq)
			req = &protos.GrpcReq{
				Data: any,
			}
			actualResp, err := us.ResetPassword(ctx, req)
			So(err, ShouldBeNil)
			So(actualResp.GetCode(), ShouldEqual, 200)
		})

		Convey("failed_by_newPassword_same_as_old", func() {
			resetPwdReq := &protos.ResetPasswordReq{
				OldPassword:      "1234567890",
				NewPassword:      "1234567890",
				VerificationCode: "112233",
				TargetAccount: &protos.ResetPasswordReq_Email{
					Email: "123@qq.com",
				},
			}
			any, _ := utils.MarshalMessageToAny(resetPwdReq)
			req = &protos.GrpcReq{
				Data: any,
			}
			actualResp, err := us.ResetPassword(ctx, req)
			So(err, ShouldBeNil)
			So(actualResp.GetMessage(), ShouldEqual, "the new password cannot be the same as the old one")
		})
	})
}
