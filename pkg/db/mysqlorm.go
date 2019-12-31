package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	config "goim-pro/configs"
	"goim-pro/pkg/logs"
	"strconv"
	"sync"
)

// connection pool
type MysqlConnectionPool struct{}

var mysqlInstance *MysqlConnectionPool

// use `sync.Once` aim to control the instance will only use once on multiple thread environment
var mysqlOnce sync.Once

var (
	db     *gorm.DB
	dbErr  error
	logger = logs.GetLogger("PANIC")
)

func init() {
	if err := initConnectionPool(); err != nil {
		logger.Panicf("initial mysql connection pool error: %v\n", err)
	}
}

/* the method to init mysql connection pool */
func initConnectionPool() error {
	var (
		dbUserName        = config.GetMysqlDbUserName()
		dbPassword        = config.GetMysqlDbPassword()
		dbUri             = config.GetMysqlDbUri()
		dbPort            = config.GetMysqlDbPort()
		dbName            = config.GetMysqlDbName()
		dbEngine          = config.GetMysqlDbEngine()
		dbMaxIdleConns, _ = strconv.Atoi(config.GetMysqlDbMaxIdleConns())
		dbMaxOpenConns, _ = strconv.Atoi(config.GetMysqlDbMaxOpenConns())
		dbEnableLogMode   = config.GetMysqlDbEnableLogMode()
	)
	connUrl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", dbUserName, dbPassword, dbUri, dbPort, dbName)
	db, dbErr = gorm.Open("mysql", connUrl)
	if dbErr != nil {
		logger.Errorf("mysql connect fail: %v\n", dbErr)
		return dbErr
	}
	logger.Infof("mysql connect successful: %s\n", connUrl)

	engine := fmt.Sprintf("ENGINE=%s", dbEngine)
	db.Set("gorm:table_options", engine)
	db.DB().SetMaxIdleConns(dbMaxIdleConns)
	db.DB().SetMaxOpenConns(dbMaxOpenConns)

	db.LogMode(dbEnableLogMode)

	// 关闭数据库，db会被多个goroutine共享，可以不调用
	// defer db.Close()
	return nil
}

/* to get mysql connect from pool as single case */
func GetMysqlConnection() *MysqlConnectionPool {
	mysqlOnce.Do(func() {
		mysqlInstance = &MysqlConnectionPool{}
	})
	return mysqlInstance
}

func (m *MysqlConnectionPool) GetMysqlDBInstance() *gorm.DB {
	return db
}
