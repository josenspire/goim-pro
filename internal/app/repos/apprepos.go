package repos

import (
	"goim-pro/internal/app/repos/address"
	"goim-pro/internal/app/repos/contact"
	"goim-pro/internal/app/repos/user"
	mysqlsrv "goim-pro/pkg/db/mysql"
)

type RepoServer struct {
	UserRepo    user.IUserRepo
	AddressRepo address.IAddress
	ContactRepo contact.IContactRepo
}

func New() *RepoServer {
	mysqlDB := mysqlsrv.NewMysqlConnection().GetMysqlInstance()
	return &RepoServer{
		UserRepo:    user.NewUserRepo(mysqlDB),
		AddressRepo: address.New(),
		ContactRepo: contact.NewContactRepo(mysqlDB),
	}
}
