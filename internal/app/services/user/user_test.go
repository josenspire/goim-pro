package usersrv

import (
	"context"
	protos "goim-pro/api/protos/salty"
	"goim-pro/config"
	"goim-pro/internal/app/models"
	"goim-pro/internal/app/repos/user"
	redsrv "goim-pro/pkg/db/redis"
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

var modelUserProfile1 = &models.UserProfile{
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
var modelUserProfile2 = &models.UserProfile{
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

var modelUser1 = &models.User{
	UserProfile: *modelUserProfile1,
}
var modelUser2 = &models.User{
	UserProfile: *modelUserProfile2,
}

func Test_Register(t *testing.T) {
	m := &MockUserRepo{}
	m.On("IsTelephoneOrEmailRegistered", "13631210001", "123@qq.com").Return(true, nil)
	m.On("IsTelephoneOrEmailRegistered", "13631210002", "12345@qq.com").Return(false, nil)
	m.On("Register", &models.User{
		Password:    "1234567890",
		UserProfile: *modelUserProfile1,
	}).Return(nil)
	m.On("Register", &models.User{
		Password:    "1234567890",
		UserProfile: *modelUserProfile2,
	}).Return(nil)

	redisDB := redsrv.NewRedisConnection()
	_ = redisDB.Connect()

	mr := &redsrv.MockCmdable{}
	mr.On("Get", "0-13631210001").Return(123456)
	mr.On("Gel", "0-13631210002").Return(123456)

	mr.On("Del", "0-13631210001").Return(1)
	mr.On("Del", "0-13631210002").Return(1)

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
	errEnPassword, _ := crypto.AESEncrypt("1234567891", config.GetApiSecretKey())
	newEnPassword, _ := crypto.AESEncrypt("1122334455", config.GetApiSecretKey())

	m := &MockUserRepo{}
	m.On("QueryByTelephoneAndPassword", "13631210001", enPassword).Return(modelUser1, nil)
	m.On("QueryByEmailAndPassword", "123@qq.com", enPassword).Return(modelUser1, nil)
	m.On("QueryByTelephoneAndPassword", "13631210001", errEnPassword).Return(&user.UserImpl{}, utils.ErrAccountOrPwdInvalid)

	m.On("IsTelephoneOrEmailRegistered", "13631210001", "").Return(true, nil)
	m.On("IsTelephoneOrEmailRegistered", "", "123@qq.com").Return(true, nil)
	m.On("IsTelephoneOrEmailRegistered", "", "123456@qq.com").Return(false, utils.ErrUserNotExists)

	m.On("ResetPasswordByTelephone", "13631210001", newEnPassword).Return(nil)
	m.On("ResetPasswordByEmail", "123@qq.com", newEnPassword).Return(nil)

	Convey("testing_grpc_user_reset_password", t, func() {
		var ctx context.Context
		var req *protos.GrpcReq
		us := &userService{
			userRepo: m,
		}
		Convey("user_reset_password_successful_by_telephone_with_old_password", func() {
			resetPwdReq := &protos.ResetPasswordReq{
				NewPassword: "1122334455",
				ResetCertificate: &protos.ResetPasswordReq_OldPassword{
					OldPassword: "1234567890",
				},
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
		Convey("user_reset_password_successful_by_email_with_verification_code", func() {
			resetPwdReq := &protos.ResetPasswordReq{
				NewPassword: "1122334455",
				ResetCertificate: &protos.ResetPasswordReq_VerificationCode{
					VerificationCode: "112233",
				},
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
				NewPassword: "1234567890",
				ResetCertificate: &protos.ResetPasswordReq_OldPassword{
					OldPassword: "1234567890",
				},
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
		Convey("failed_by_not_exist_account", func() {
			resetPwdReq := &protos.ResetPasswordReq{
				NewPassword: "1234567890",
				ResetCertificate: &protos.ResetPasswordReq_OldPassword{
					OldPassword: "1234567891",
				},
				TargetAccount: &protos.ResetPasswordReq_Email{
					Email: "123456@qq.com",
				},
			}
			any, _ := utils.MarshalMessageToAny(resetPwdReq)
			req = &protos.GrpcReq{
				Data: any,
			}
			actualResp, _ := us.ResetPassword(ctx, req)
			So(actualResp.GetMessage(), ShouldEqual, utils.ErrUserNotExists.Error())
		})
		Convey("failed_by_invalid_oldPassword", func() {
			resetPwdReq := &protos.ResetPasswordReq{
				NewPassword: "1234567890",
				ResetCertificate: &protos.ResetPasswordReq_OldPassword{
					OldPassword: "1234567891",
				},
				TargetAccount: &protos.ResetPasswordReq_Telephone{
					Telephone: "13631210001",
				},
			}
			any, _ := utils.MarshalMessageToAny(resetPwdReq)
			req = &protos.GrpcReq{
				Data: any,
			}
			actualResp, _ := us.ResetPassword(ctx, req)
			So(actualResp.GetMessage(), ShouldEqual, utils.ErrAccountOrPwdInvalid.Error())
		})
	})
}

func Test_userService_GetUserInfo(t *testing.T) {
	m := &MockUserRepo{}
	m.On("FindByUserId", "13631210001").Return(modelUser1, nil)

	Convey("testing_grpc_get_user_by_userId", t, func() {
		var ctx context.Context
		var req *protos.GrpcReq
		us := &userService{
			userRepo: m,
		}
		Convey("should_return_user_when_given_correct_userId", func() {
			userReq := &protos.GetUserInfoReq{
				UserId: "13631210001",
			}
			any, _ := utils.MarshalMessageToAny(userReq)
			req = &protos.GrpcReq{
				Data: any,
			}
			actualResp, err := us.GetUserInfo(ctx, req)

			getUserProfileResp := &protos.GetUserInfoResp{}
			err = utils.UnMarshalAnyToMessage(actualResp.GetData(), getUserProfileResp)

			So(err, ShouldBeNil)
			So(getUserProfileResp.Profile.GetTelephone(), ShouldEqual, "13631210001")
		})
	})
}

func Test_userService_QueryUserInfo(t *testing.T) {
	userCriteria1 := &user.UserImpl{}
	userCriteria1.Telephone = "13631210001"

	userCriteria2 := &user.UserImpl{}
	userCriteria2.Email = "123@qq.com"

	userCriteria3 := &user.UserImpl{}
	userCriteria3.Telephone = "13631210012"

	m := &MockUserRepo{}
	m.On("FindOneUser", userCriteria1).Return(modelUser1, nil)
	m.On("FindOneUser", userCriteria2).Return(modelUser1, nil)
	m.On("FindOneUser", userCriteria3).Return(&user.UserImpl{}, utils.ErrUserNotExists)

	Convey("testing_grpc_query_user_info", t, func() {
		var ctx context.Context
		var req *protos.GrpcReq
		us := &userService{
			userRepo: m,
		}
		Convey("should_return_user_when_given_correct_telephone", func() {
			userReq := &protos.QueryUserInfoReq{
				TargetAccount: &protos.QueryUserInfoReq_Telephone{
					Telephone: "13631210001",
				},
			}
			any, _ := utils.MarshalMessageToAny(userReq)
			req = &protos.GrpcReq{
				Data: any,
			}
			actualResp, err := us.QueryUserInfo(ctx, req)

			queryUserInfoResp := &protos.QueryUserInfoResp{}
			err = utils.UnMarshalAnyToMessage(actualResp.GetData(), queryUserInfoResp)

			So(err, ShouldBeNil)
			So(queryUserInfoResp.Profile.GetTelephone(), ShouldEqual, "13631210001")
		})
		Convey("should_return_user_when_given_correct_email", func() {
			userReq := &protos.QueryUserInfoReq{
				TargetAccount: &protos.QueryUserInfoReq_Email{
					Email: "123@qq.com",
				},
			}
			any, _ := utils.MarshalMessageToAny(userReq)
			req = &protos.GrpcReq{
				Data: any,
			}
			actualResp, err := us.QueryUserInfo(ctx, req)

			queryUserInfoResp := &protos.QueryUserInfoResp{}
			err = utils.UnMarshalAnyToMessage(actualResp.GetData(), queryUserInfoResp)

			So(err, ShouldBeNil)
			So(queryUserInfoResp.Profile.GetTelephone(), ShouldEqual, "13631210001")
		})
		Convey("should_return_err_when_given_un_exists_telephone", func() {
			userReq := &protos.QueryUserInfoReq{
				TargetAccount: &protos.QueryUserInfoReq_Telephone{
					Telephone: "13631210012",
				},
			}
			any, _ := utils.MarshalMessageToAny(userReq)
			req = &protos.GrpcReq{
				Data: any,
			}
			actualResp, _ := us.QueryUserInfo(ctx, req)

			queryUserInfoResp := &protos.QueryUserInfoResp{}
			_ = utils.UnMarshalAnyToMessage(actualResp.GetData(), queryUserInfoResp)

			So(actualResp.GetMessage(), ShouldEqual, utils.ErrUserNotExists.Error())
		})
	})
}

func Test_userService_UpdateUserInfo(t *testing.T) {
	criteria1 := &models.User{}
	criteria1.UserId = "13631210001"

	criteria2 := &models.User{}
	criteria2.UserId = "13631210002"

	newProfile1 := models.UserProfile{
		UserId:      "13631210001",
		Telephone:   "13631214444",
		Email:       "123456@qq.com",
		Nickname:    "JAMESYANG01",
		Avatar:      "",
		Description: "",
		Sex:         "0",
		Birthday:    0,
		Location:    "ZHA",
	}
	newProfile2 := models.UserProfile{
		UserId:      "13631210002",
		Telephone:   "13631214444",
		Email:       "123456@qq.com",
		Nickname:    "JAMESYANG01",
		Avatar:      "",
		Description: "",
		Sex:         "0",
		Birthday:    0,
		Location:    "ZHA",
	}

	m := &MockUserRepo{}
	m.On("FindByUserId", "13631210001").Return(modelUser1, nil)
	m.On("FindByUserId", "13631210002").Return(&user.UserImpl{}, nil)
	m.On("FindOneAndUpdateProfile", criteria1, utils.TransformStructToMap(newProfile1)).Return(nil)
	m.On("FindOneAndUpdateProfile", criteria2, utils.TransformStructToMap(newProfile2)).Return(nil)

	Convey("testing_grpc_update_user_profile", t, func() {
		var ctx context.Context
		var req *protos.GrpcReq
		us := &userService{
			userRepo: m,
		}
		Convey("should_update_profile_success_when_given_correct_newProfile", func() {
			newProfile1 := &protos.UserProfile{
				UserId:      "13631210001",
				Telephone:   "13631214444",
				Email:       "123456@qq.com",
				Nickname:    "JAMESYANG01",
				Avatar:      "",
				Description: "",
				Sex:         0,
				Birthday:    0,
				Location:    "ZHA",
			}
			updateUserReq := &protos.UpdateUserInfoReq{
				Profile: newProfile1,
			}
			any, _ := utils.MarshalMessageToAny(updateUserReq)
			req = &protos.GrpcReq{
				Data: any,
			}
			actualResp, err := us.UpdateUserInfo(ctx, req)

			So(err, ShouldBeNil)
			So(actualResp.Code, ShouldEqual, 200)
		})

		Convey("should_update_profile_failed_when_given_incorrect_profile_with_userId", func() {
			newProfile2 := &protos.UserProfile{
				UserId:      "13631210003",
				Telephone:   "13631214444",
				Email:       "123456@qq.com",
				Nickname:    "JAMESYANG01",
				Avatar:      "",
				Description: "",
				Sex:         0,
				Birthday:    0,
				Location:    "ZHA",
			}
			updateUserReq := &protos.UpdateUserInfoReq{
				Profile: newProfile2,
			}
			any, _ := utils.MarshalMessageToAny(updateUserReq)
			req = &protos.GrpcReq{
				Data: any,
			}
			actualResp, _ := us.UpdateUserInfo(ctx, req)

			So(actualResp.Code, ShouldEqual, 400)
			So(actualResp.Message, ShouldEqual, utils.ErrInvalidUserId)
		})
	})
}
