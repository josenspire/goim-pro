package services

import (
	"goim-pro/api/protos"
	demowaiter "goim-pro/internal/app/services/demowaiter"
	usersrv "goim-pro/internal/app/services/user"
)

type Service struct {
	UserServer   protos.UserServiceServer
	WaiterServer protos.WaiterServer
}

func NewService() *Service {
	return &Service{
		UserServer:   usersrv.New(),
		WaiterServer: demowaiter.New(),
	}
}
