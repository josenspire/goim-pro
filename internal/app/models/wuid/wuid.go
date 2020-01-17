package wuid

import (
	"database/sql"
	"github.com/edwingeng/wuid/mysql/wuid"
	mysqlsrv "goim-pro/pkg/db/mysql"
	"goim-pro/pkg/logs"
	"strconv"
)

type Wuid struct {
	H uint  `gorm:"primary_key; AUTO_INCREMENT; not null;"`
	X uint8 `gorm:"unique_index; default: 0; not null;"`
}

var logger = logs.GetLogger("ERROR")
var g *wuid.WUID

func init() {
	wuid.WithSection(5)
	g = wuid.NewWUID("default", nil)
}

func NewWUID() string {
	newDB := func() (*sql.DB, bool, error) {
		db := mysqlsrv.NewMysqlConnection().GetMysqlInstance().DB()
		return db, false, nil
	}
	// setup
	if err := g.LoadH28FromMysql(newDB, "wuids"); err != nil {
		logger.Errorf("load wuid error: %s", err.Error())
		return "0"
	}
	// generate
	return strconv.Itoa(int(g.Next()))
}
