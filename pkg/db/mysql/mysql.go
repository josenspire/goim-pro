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
	mysqlDB       *gorm.DB
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
	mysqlDB, err = gorm.Open("mysql", connUrl)
	if err != nil {
		return
	}
	logger.Infof("[mysql] connect successful: %s", connUrl)

	engine := fmt.Sprintf("ENGINE=%s", dbEngine)
	mysqlDB.Set("gorm:table_options", engine)
	mysqlDB.DB().SetMaxIdleConns(dbMaxIdleConns)
	mysqlDB.DB().SetMaxOpenConns(dbMaxOpenConns)

	mysqlDB.LogMode(dbEnableLogMode)

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
	err = initConnectionPool()
	return
}

func (m *MysqlConnectionPool) GetMysqlInstance() *gorm.DB {
	return mysqlDB
}
