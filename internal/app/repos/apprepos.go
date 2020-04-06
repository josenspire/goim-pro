package repos

import (
	"github.com/jinzhu/gorm"
	"goim-pro/internal/app/repos/address"
	"goim-pro/internal/app/repos/contact"
	"goim-pro/internal/app/repos/group"
	"goim-pro/internal/app/repos/user"
)

type RepoServer struct {
	UserRepo    user.IUserRepo
	AddressRepo address.IAddress
	ContactRepo contact.IContactRepo
	GroupRepo   group.IGroupRepo
}

func New(mysqlDB *gorm.DB) *RepoServer {
	return &RepoServer{
		UserRepo:    user.NewUserRepo(mysqlDB),
		AddressRepo: address.New(),
		ContactRepo: contact.NewContactRepo(mysqlDB),
		GroupRepo:   group.NewGroupRepo(mysqlDB),
	}
}
