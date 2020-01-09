package redsrv

import (
	"testing"
)

func TestRedisServiceConnection_GetRedisClient(t *testing.T) {
	redisDB := NewRedisConnection()
	err := redisDB.Connect()
	if err != nil {
		t.FailNow()
	} else {
		db := redisDB.GetRedisClient()
		if db == nil {
			t.FailNow()
		}
	}
}

