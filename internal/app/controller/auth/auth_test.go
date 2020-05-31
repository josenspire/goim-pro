package auth

import (
	"context"
	. "github.com/smartystreets/goconvey/convey"
	protos "goim-pro/api/protos/salty"
	cserr "goim-pro/internal/app/models/errors"
	errmsg "goim-pro/pkg/errors"
	"goim-pro/pkg/utils"
	"testing"
)

func Test_authServer_ObtainTelephoneSMSCode(t *testing.T) {
	telephone1 := "13631210000"
	telephone2 := "13631210001"
	telephone3 := "13631210002"
	registerType := protos.SMSOperationType_REGISTER
	resetPasswordType := protos.SMSOperationType_RESET_PASSWORD

	existsErr := &cserr.TError{
		Code:   protos.StatusCode_STATUS_ACCOUNT_EXISTS,
		Detail: errmsg.ErrTelephoneExists.Error(),
	}
	notExistsErr := &cserr.TError{
		Code:   protos.StatusCode_STATUS_ACCOUNT_NOT_EXISTS,
		Detail: errmsg.ErrAccountNotExists.Error(),
	}

	authController := New()
	m := &MockAuthService{}
	m.On("ObtainSMSCode", telephone1, registerType).Return("123401", nil)
	m.On("ObtainSMSCode", telephone2, registerType).Return("", existsErr)
	m.On("ObtainSMSCode", telephone3, resetPasswordType).Return("", notExistsErr)

	authService = m

	Convey("Testing_ObtainTelephoneSMSCode", t, func() {
		ctx := context.Background()
		Convey("should_call_and_return_register_code_when_given_unExists_telephone_account", func() {
			authReq := &protos.ObtainTelephoneSMSCodeReq{
				OperationType: protos.SMSOperationType_REGISTER,
				Telephone:     "13631210000",
			}
			anyData, _ := utils.MarshalMessageToAny(authReq)
			grpcReq := &protos.GrpcReq{
				DeviceId: "XIAOMI 10",
				Token:    "1234657890",
				Data:     anyData,
			}
			resp, grpcErr := authController.ObtainTelephoneSMSCode(ctx, grpcReq)
			So(grpcErr, ShouldBeNil)
			So(resp.Code, ShouldEqual, protos.StatusCode_STATUS_OK)
		})
		Convey("should_call_and_return_error_when_given_exists_telephone_to_register", func() {
			authReq := &protos.ObtainTelephoneSMSCodeReq{
				OperationType: protos.SMSOperationType_REGISTER,
				Telephone:     "13631210001",
			}
			anyData, _ := utils.MarshalMessageToAny(authReq)
			grpcReq := &protos.GrpcReq{
				DeviceId: "XIAOMI 10",
				Token:    "1234657890",
				Data:     anyData,
			}
			resp, grpcErr := authController.ObtainTelephoneSMSCode(ctx, grpcReq)
			So(grpcErr, ShouldBeNil)
			So(resp.Code, ShouldEqual, protos.StatusCode_STATUS_ACCOUNT_EXISTS)
			So(resp.Message, ShouldEqual, errmsg.ErrTelephoneExists.Error())
		})
		Convey("should_return_err_when_given_unExists_telephone_to_reset_password", func() {
			authReq := &protos.ObtainTelephoneSMSCodeReq{
				OperationType: protos.SMSOperationType_RESET_PASSWORD,
				Telephone:     "13631210002",
			}
			anyData, _ := utils.MarshalMessageToAny(authReq)
			grpcReq := &protos.GrpcReq{
				DeviceId: "XIAOMI 10",
				Token:    "1234657890",
				Data:     anyData,
			}
			resp, grpcErr := authController.ObtainTelephoneSMSCode(ctx, grpcReq)
			So(grpcErr, ShouldBeNil)
			So(resp.Code, ShouldEqual, protos.StatusCode_STATUS_ACCOUNT_NOT_EXISTS)
			So(resp.Message, ShouldEqual, errmsg.ErrAccountNotExists.Error())
		})
	})
}

func Test_authServer_VerifyTelephoneSMSCode(t *testing.T) {
	telephone1 := "13631210001"
	registerType := protos.SMSOperationType_REGISTER
	registerCode := "123401"
	invalidRegisterCode := "123402"

	tErr := &cserr.TError{
		Code:   protos.StatusCode_STATUS_BAD_REQUEST,
		Detail: errmsg.ErrInvalidVerificationCode.Error(),
	}

	m := &MockAuthService{}
	m.On("VerifySMSCode", telephone1, registerType, registerCode).Return(true, nil)
	m.On("VerifySMSCode", telephone1, registerType, invalidRegisterCode).Return(false, tErr)
	authService = m

	authController := new(authServer)

	Convey("Testing_VerifyTelephoneSMSCode", t, func() {
		ctx := context.Background()
		Convey("should_verify_successful_when_given_correct_telephone_operationType_and_code", func() {
			authReq := &protos.VerifyTelephoneSMSCodeReq{
				OperationType: protos.SMSOperationType_REGISTER,
				Telephone:     "13631210001",
				SmsCode:       "123401",
			}
			anyData, _ := utils.MarshalMessageToAny(authReq)
			grpcReq := &protos.GrpcReq{
				DeviceId: "XIAOMI 10",
				Token:    "1234657890",
				Data:     anyData,
			}
			resp, tErr := authController.VerifyTelephoneSMSCode(ctx, grpcReq)
			So(tErr, ShouldBeNil)
			So(resp.Code, ShouldEqual, protos.StatusCode_STATUS_OK)
		})
		Convey("should_verify_fail_when_given_incorrect_sms_code", func() {
			authReq := &protos.VerifyTelephoneSMSCodeReq{
				OperationType: protos.SMSOperationType_REGISTER,
				Telephone:     "13631210001",
				SmsCode:       "123402",
			}
			anyData, _ := utils.MarshalMessageToAny(authReq)
			grpcReq := &protos.GrpcReq{
				DeviceId: "XIAOMI 10",
				Token:    "1234657890",
				Data:     anyData,
			}
			resp, grpcErr := authController.VerifyTelephoneSMSCode(ctx, grpcReq)
			So(grpcErr, ShouldBeNil)
			So(resp.Code, ShouldEqual, protos.StatusCode_STATUS_BAD_REQUEST)
			So(resp.Message, ShouldEqual, errmsg.ErrInvalidVerificationCode.Error())
		})
	})
}
