package mysqlsrv

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"goim-pro/config"
	"strconv"
)

type BaseMysql struct {
	gorm.DB
}

func newBaseMysql() *gorm.DB {
	var (
		dbUserName        = config.GetMysqlDBUserName()
		dbPassword        = config.GetMysqlDBPassword()
		dbUri             = config.GetMysqlDBUri()
		dbPort            = config.GetMysqlDBPort()
		dbName            = config.GetMysqlDBName()
		dbEngine          = config.GetMysqlDBEngine()
		dbMaxIdleConns, _ = strconv.Atoi(config.GetMysqlDBMaxIdleConns())
		dbMaxOpenConns, _ = strconv.Atoi(config.GetMysqlDBMaxOpenConns())
		dbEnableLogMode   = config.GetMysqlDBEnableLogMode()
	)
	connUrl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", dbUserName, dbPassword, dbUri, dbPort, dbName)
	mysqlDB, err := gorm.Open("mysql", connUrl)
	if err != nil {
		return nil
	}
	logger.Infof("[mysql] connect successful: %s", connUrl)

	engine := fmt.Sprintf("ENGINE=%s", dbEngine)
	mysqlDB.Set("gorm:table_options", engine)
	mysqlDB.DB().SetMaxIdleConns(dbMaxIdleConns)
	mysqlDB.DB().SetMaxOpenConns(dbMaxOpenConns)

	mysqlDB.LogMode(dbEnableLogMode)

	return mysqlDB
}

func (m *BaseMysql) FindOne(condition interface{}, opts ...interface{}) (inter interface{}, err error) {
	return
}
