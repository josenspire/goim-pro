package mysqlsrv

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"goim-pro/config"
	"goim-pro/pkg/logs"
	"strconv"
)

var (
	logger = logs.GetLogger("INFO")

	mysqlDB *gorm.DB
)

/* the method to init mysql connection pool */
func NewMysql() {
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
	mysqlDB = newBaseMysql(options)
}

/* to get mysql connect from pool as single case */
func GetMysql() *gorm.DB {
	return mysqlDB
}
