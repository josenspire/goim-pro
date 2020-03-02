package repos

import (
	"goim-pro/internal/app/repos/address"
	"goim-pro/internal/app/repos/user"
	mysqlsrv "goim-pro/pkg/db/mysql"
)

type RepoServer struct {
	UserRepo    user.IUserRepo
	AddressRepo address.IAddress
}

func New() *RepoServer {
	mysqlDB := mysqlsrv.NewMysqlConnection().GetMysqlInstance()
	return &RepoServer{
		UserRepo:    user.NewUserRepo(mysqlDB),
		AddressRepo: address.New(),
	}
}
