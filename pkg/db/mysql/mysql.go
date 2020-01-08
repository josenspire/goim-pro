package mysqlsrv

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	config "goim-pro/config"
	"goim-pro/pkg/logs"
	"strconv"
	"sync"
)

// connection pool
type MysqlConnectionPool struct{}

var (
	// use `sync.Once` aim to control the instance will only use once on multiple thread environment
	mysqlOnce     sync.Once
	db            *gorm.DB
	mysqlInstance *MysqlConnectionPool
	logger        = logs.GetLogger("INFO")
)

/* the method to init mysql connection pool */
func initConnectionPool() (err error) {
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
	db, err = gorm.Open("mysql", connUrl)
	if err != nil {
		logger.Errorf("[mysql] connect fail: %v\n", err)
		return
	}
	logger.Infof("[mysql] connect successful: %s\n", connUrl)

	engine := fmt.Sprintf("ENGINE=%s", dbEngine)
	db.Set("gorm:table_options", engine)
	db.DB().SetMaxIdleConns(dbMaxIdleConns)
	db.DB().SetMaxOpenConns(dbMaxOpenConns)

	db.LogMode(dbEnableLogMode)

	// 关闭数据库，db会被多个goroutine共享，可以不调用
	// defer db.Close()
	return
}

/* to get mysql connect from pool as single case */
func NewMysqlConnection() *MysqlConnectionPool {
	mysqlOnce.Do(func() {
		mysqlInstance = &MysqlConnectionPool{}
	})
	return mysqlInstance
}

func (m *MysqlConnectionPool) Connect() (err error) {
	if err := initConnectionPool(); err != nil {
		logger.Panicf("initial mysql connection pool error: %v\n", err)
	}
	return
}

func GetMysqlInstance() *gorm.DB {
	return db
}
