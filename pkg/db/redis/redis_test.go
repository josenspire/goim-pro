package redsrv

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestRedisServiceConnection_GetRedisClient(t *testing.T) {
	redisDB := NewRedis()

	Convey("TestRedis", t, func() {
		Convey("should_set_and_get_data", func() {
			err := redisDB.RSet("TEST_01", "ASDFGHJKL", time.Duration(2)*time.Second)
			So(err, ShouldBeNil)
			value := redisDB.RGet("TEST_01")
			So(value, ShouldEqual, "ASDFGHJKL")

			redisDB.RDel("TEST_01")
		})
	})
}
