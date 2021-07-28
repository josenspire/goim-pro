package redsrv

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func init() {
	NewRedis()
}

func TestRedisServiceConnection_GetRedisClient(t *testing.T) {
	myRedis := GetRedis()

	Convey("TestRedis", t, func() {
		Convey("should_set_and_get_data", func() {
			err := myRedis.RSet("TEST_01", "ASDFGHJKL", time.Duration(2)*time.Second)
			So(err, ShouldBeNil)
			value := myRedis.RGet("TEST_01")
			So(value, ShouldEqual, "ASDFGHJKL")

			myRedis.RDel("TEST_01")
		})
	})
}
