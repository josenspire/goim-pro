package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"sync"
)

// connection pool
type MysqlConnectionPool struct{}

var mysqlInstance *MysqlConnectionPool

// use `sync.Once` aim to control the instance will only use once on multiple thread environment
var mysqlOnce sync.Once

var (
	db    *gorm.DB
	dbErr error
)

// TODO: should be refactor to read from config file
const (
	dbUserName      = "root"
	dbPassword      = "Password1!"
	dbUri           = "127.0.0.1"
	dbPort          = "3306"
	dbName          = "goim"
	dbEngine        = "InnoDB"
	dbMaxIdleConns  = 10
	dbMaxOpenConns  = 30
	dbEnableLogMode = true
)

/* to get mysql connect from pool as single case */
func GetMysqlConnection() *MysqlConnectionPool {
	mysqlOnce.Do(func() {
		mysqlInstance = &MysqlConnectionPool{}
	})
	return mysqlInstance
}

/* the method to init mysql connection pool */
func (m *MysqlConnectionPool) InitConnectionPool() (bool, error) {
	connUrl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", dbUserName, dbPassword, dbUri, dbPort, dbName)
	db, dbErr = gorm.Open("mysql", connUrl)
	if dbErr != nil {
		fmt.Printf("mysql connect fail: %v\n", dbErr)
		return false, dbErr
	}
	fmt.Printf("mysql connect successful: %s\n", connUrl)

	engine := fmt.Sprintf("ENGINE=%s", dbEngine)
	db.Set("gorm:table_options", engine)
	db.DB().SetMaxIdleConns(dbMaxIdleConns)
	db.DB().SetMaxOpenConns(dbMaxOpenConns)

	db.LogMode(dbEnableLogMode)

	// 关闭数据库，db会被多个goroutine共享，可以不调用
	// defer db.Close()
	return true, nil
}

func (m *MysqlConnectionPool) GetMysqlDBInstance() *gorm.DB {
	return db
}
