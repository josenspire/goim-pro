package mysqlsrv

import (
	"testing"
)

func TestMysqlConnectionPool_GetMysqlInstance(t *testing.T) {
	err := NewMysqlConnection().Connect()
	if err != nil {
		t.FailNow()
	} else {
		db := GetMysqlInstance()
		if db == nil {
			t.FailNow()
		}
	}
}
