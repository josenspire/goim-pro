package db

import "testing"

func TestMysqlConnectionPool_GetMysqlDBInstance(t *testing.T) {
	err := GetMysqlConnection().InitConnectionPool()
	if err != nil {
		t.FailNow()
	}
	db := GetMysqlConnection().GetMysqlDBInstance()
	if db == nil {
		t.FailNow()
	}
}
