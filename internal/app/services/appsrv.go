package services

import (
	example "goim-pro/api/protos/example"
	protos "goim-pro/api/protos/saltyv2"
	waitersrv "goim-pro/internal/app/services/demowaiter"
	usersrv "goim-pro/internal/app/services/user"
)

type Service struct {
	UserServer   protos.UserServiceServer
	WaiterServer example.WaiterServer
}

func NewService() *Service {
	return &Service{
		WaiterServer: waitersrv.New(),
		UserServer:   usersrv.New(),
	}
}
