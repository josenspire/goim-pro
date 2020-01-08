package redsrv

import (
	"testing"
)

func TestRedisServiceConnection_GetRedisClient(t *testing.T) {
	err := NewRedisService().Connect()
	if err != nil {
		t.FailNow()
	} else {
		db := GetRedisClient()
		if db == nil {
			t.FailNow()
		}
	}
}

