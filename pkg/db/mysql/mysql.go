package mysqlsrv

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"goim-pro/config"
	"goim-pro/pkg/logs"
	"strconv"
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
	dbMaxIdleConns, _ := strconv.Atoi(config.GetMysqlDBMaxIdleConns())
	dbMaxOpenConns, _ := strconv.Atoi(config.GetMysqlDBMaxOpenConns())

	var options *mysqlOptions = &mysqlOptions{
		dbUserName:      config.GetMysqlDBUserName(),
		dbPassword:      config.GetMysqlDBPassword(),
		dbUri:           config.GetMysqlDBUri(),
		dbPort:          config.GetMysqlDBPort(),
		dbName:          config.GetMysqlDBName(),
		dbEngine:        config.GetMysqlDBEngine(),
		dbMaxIdleConns:  dbMaxIdleConns,
		dbMaxOpenConns:  dbMaxOpenConns,
		dbEnableLogMode: config.GetMysqlDBEnableLogMode(),
	}
	return newBaseMysql(options)
}
