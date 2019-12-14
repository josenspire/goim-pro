package db

import "testing"

func TestMysqlConnectionPool_GetMysqlDBInstance(t *testing.T) {
	initResult, _ := GetMysqlConnection().InitConnectionPool()
	if !initResult {
		t.FailNow()
	}
	db := GetMysqlConnection().GetMysqlDBInstance()
	if db == nil {
		t.FailNow()
	}
}
