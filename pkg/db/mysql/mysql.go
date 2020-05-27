package mysqlsrv

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"goim-pro/pkg/logs"
	"sync"
)

var (
	logger = logs.GetLogger("INFO")

	// use `sync.Once` aim to control the instance will only use once on multiple thread environment
	mysqlOnce sync.Once
	mysqlDB   *gorm.DB
)

/* to get mysql connect from pool as single case */
func NewMysql() *gorm.DB {
	mysqlOnce.Do(func() {
		mysqlDB = connect()
	})
	return mysqlDB
}

/* the method to init mysql connection pool */
func connect() *gorm.DB {
	return newBaseMysql()
}

//func (m *MysqlConnectionPool) Connect() (err error) {
//	err = initConnectionPool()
//	return
//}
//
//func (m *MysqlConnectionPool) GetMysqlInstance() *gorm.DB {
//	return mysqlDB
//}
