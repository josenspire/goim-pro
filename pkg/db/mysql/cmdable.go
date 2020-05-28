package mysqlsrv

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type BaseMysql struct {
	gorm.DB
}

type mysqlOptions struct {
	dbUserName      string
	dbPassword      string
	dbUri           string
	dbPort          string
	dbName          string
	dbEngine        string
	dbMaxIdleConns  int
	dbMaxOpenConns  int
	dbEnableLogMode bool
}

func newBaseMysql(options *mysqlOptions) *gorm.DB {
	connUrl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", options.dbUserName, options.dbPassword, options.dbUri, options.dbPort, options.dbName)
	mysqlDB, err := gorm.Open("mysql", connUrl)
	if err != nil {
		return nil
	}
	logger.Infof("[mysql] connect successful: %s", connUrl)

	engine := fmt.Sprintf("ENGINE=%s", options.dbEngine)
	mysqlDB.Set("gorm:table_options", engine)
	mysqlDB.DB().SetMaxIdleConns(options.dbMaxIdleConns)
	mysqlDB.DB().SetMaxOpenConns(options.dbMaxOpenConns)

	mysqlDB.LogMode(options.dbEnableLogMode)

	return mysqlDB
}
