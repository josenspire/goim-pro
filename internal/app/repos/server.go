package repos

import (
	"github.com/jinzhu/gorm"
	"goim-pro/internal/app/repos/address"
	"goim-pro/internal/app/repos/user"
	"goim-pro/pkg/db"
	"goim-pro/pkg/logs"
)

var mysqlDB *gorm.DB
var logger = logs.GetLogger("ERROR")

type RepoServer struct {
	UserRepo    user.IUserRepo
	AddressRepo address.IAddress
}

func init() {
	mysqlDB = db.GetMysqlConnection().GetMysqlDBInstance()
	if err := initialMysqlTables(mysqlDB); err != nil {
		panic(err)
	}
}

func New() *RepoServer {
	return &RepoServer{
		UserRepo:    user.NewUserRepo(mysqlDB),
		AddressRepo: address.New(mysqlDB),
	}
}

func initialMysqlTables(db *gorm.DB) (err error) {
	if !db.HasTable(user.User{}) {
		err = db.Set(
			"gorm:table_options",
			"ENGINE=InnoDB DEFAULT CHARSET=utf8",
		).CreateTable(
			user.User{},
		).Error
		if err != nil {
			logger.Errorf("initial mysql tables [users] error: %v\n", err)
			return
		}
	}
	return
}
