package services

import (
	demo "goim-pro/api/protos/example"
	protos "goim-pro/api/protos/salty"
	authsrv "goim-pro/internal/app/services/auth"
	waitersrv "goim-pro/internal/app/services/demowaiter"
	usersrv "goim-pro/internal/app/services/user"
)

type Service struct {
	WaiterServer demo.WaiterServer
	SMSServer    protos.SMSServiceServer
	UserServer   protos.UserServiceServer
}

func NewService() *Service {
	return &Service{
		WaiterServer: waitersrv.New(),
		SMSServer:    authsrv.New(),
		UserServer:   usersrv.New(),
	}
}
