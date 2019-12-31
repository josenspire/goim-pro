package db

import "testing"

func TestMysqlConnectionPool_GetMysqlDBInstance(t *testing.T) {
	db := GetMysqlConnection().GetMysqlDBInstance()
	if db == nil {
		t.FailNow()
	}
}
