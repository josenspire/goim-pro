package utils

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestParseTimestampToDateTime(t *testing.T) {
	// TODO:
}

func TestParseTimestampToDateTimeStr(t *testing.T) {
	Convey("TestParseTimestampToDateTimeStr", t, func() {
		var ts int64 = 1531293019
		//var ts int64 = 1578903121862
		Convey("testing ParseTimestampToDateTimeStr", func() {
			result := ParseTimestampToDateTimeStr(ts, "2006-01-02T15:04:05+08:00")
			So(result, ShouldEqual, "2018-07-11T15:10:19+08:00")
		})
		Convey("testing ParseTimestampToDateTimeStr 222", func() {
			result := ParseTimestampToDateTimeStr(1602605419, MysqlDateTimeFormat)
			fmt.Println("-------", result)
		})
	})
}

func TestParseDateTimeStrToTimestamp(t *testing.T) {
	Convey("testing ParseDateTimeStrToTimestamp", t, func() {
		var dateTimeStr string = "2018-07-11T15:10:19+08:00"
		timestamp, err := ParseDateTimeStrToTimestamp(DefaultDateTimeFormat, dateTimeStr)
		if err != nil {
			t.FailNow()
		}
		So(timestamp, ShouldEqual, 1531293019)

		var dateTimeStr2 string = "2020-10-14T00:10:19+08:00"
		timestamp2, _ := ParseDateTimeStrToTimestamp(DefaultDateTimeFormat, dateTimeStr2)
		fmt.Println("=============", timestamp2)
	})
}
