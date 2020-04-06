package services

import (
	demo "goim-pro/api/protos/example"
	protos "goim-pro/api/protos/salty"
	authsrv "goim-pro/internal/app/services/auth"
	contactsrv "goim-pro/internal/app/services/contact"
	waitersrv "goim-pro/internal/app/services/demowaiter"
	usersrv "goim-pro/internal/app/services/user"
)

type Service struct {
	WaiterServer  demo.WaiterServer
	SMSServer     protos.SMSServiceServer
	UserServer    protos.UserServiceServer
	ContactServer protos.ContactServiceServer
	GroupServer	  protos.GroupServiceServer
}

func NewService() *Service {
	return &Service{
		WaiterServer: waitersrv.New(),
		SMSServer:    authsrv.New(),
		UserServer:   usersrv.New(),
		ContactServer: contactsrv.New(),
		GroupServer: groupsrv.New(),
	}
}
