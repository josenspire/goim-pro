package grpc

import "goim-pro/pkg/utils"

var whiteList = map[string]bool{
	"/com.salty.protos.SMSService/ObtainSMSCode": true,

	"/com.salty.protos.UserService/Register":      true,
	"/com.salty.protos.UserService/Login":         true,
	"/com.salty.protos.UserService/Logout":        true,
	"/com.salty.protos.UserService/ResetPassword": true,
}

func isOnWhiteList(fullMethodName string) (isValid bool) {
	if utils.IsEmptyStrings(fullMethodName) {
		isValid = false
	} else {
		_, isValid = whiteList[fullMethodName]
	}
	return
}
