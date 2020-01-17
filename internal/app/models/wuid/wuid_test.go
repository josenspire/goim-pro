package wuid

import (
	"fmt"
	mysqlsrv "goim-pro/pkg/db/mysql"
	"testing"
)

func TestNewWUID(t *testing.T) {
	mysqlDB := mysqlsrv.NewMysqlConnection()
	_ = mysqlDB.Connect()

	wuid1 := NewWUID()
	fmt.Println(wuid1)
	wuid2 := NewWUID()
	fmt.Println(wuid2)
}
