package redsrv

import (
	"testing"
)

func TestRedisServiceConnection_GetRedisClient(t *testing.T) {
	redisDB := NewRedisConnection()
	if redisDB.Error != nil {
		t.FailNow()
	} else {
		db := redisDB.GetRedisClient()
		if db == nil {
			t.FailNow()
		}
	}
}
