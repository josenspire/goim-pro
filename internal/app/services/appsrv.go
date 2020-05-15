package services

import (
	demo "goim-pro/api/protos/example"
	protos "goim-pro/api/protos/salty"
	"goim-pro/internal/app/controller/auth"
	"goim-pro/internal/app/controller/contact"
	"goim-pro/internal/app/controller/group"
	"goim-pro/internal/app/controller/user"
	waitersrv "goim-pro/internal/app/services/demowaiter"
)

type Service struct {
	WaiterServer  demo.WaiterServer
	SMSServer     protos.SMSServiceServer
	UserServer    protos.UserServiceServer
	ContactServer protos.ContactServiceServer
	GroupServer   protos.GroupServiceServer
}

func NewService() *Service {
	return &Service{
		WaiterServer: waitersrv.New(),
		//SMSServer:     authsrv.New(),
		//UserServer:    usersrv.New(),
		//ContactServer: contactsrv.New(),
		//GroupServer:   groupsrv.New(),
		SMSServer:     auth.New(),
		UserServer:    user.New(),
		ContactServer: contact.New(),
		GroupServer:   group.New(),
	}
}
