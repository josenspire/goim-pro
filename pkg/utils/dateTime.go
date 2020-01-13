package utils

import (
	"errors"
	"time"
)

var loc *time.Location

const (
	DefaultTimeZoneName   = "Asia/Shanghai" // default time zone as Asia/Shanghai (+08:00)
	DefaultDateTimeFormat = "2006-01-02T15:04:05+08:00"
)

func init() {
	// setup default time zone as Asia/Shanghai (+08:00)
	loc, _ = time.LoadLocation(DefaultTimeZoneName)
}

func ParseTimestampToDateTime(timestamp int64) time.Time {
	if timestamp == 0 {
		return time.Unix(0, 0)
	}
	return time.Unix(timestamp, 0)
}

// parse timestamp to date time string
// timestamp {int64}  target input timestamp
// format	{string} target output format
func ParseTimestampToDateTimeStr(timestamp int64, format string) string {
	if timestamp == 0 {
		return ""
	}
	if IsEmptyString(format) {
		format = DefaultDateTimeFormat
	}
	return time.Unix(timestamp, 0).In(loc).Format(format)
}

// parse date time string to timestamp
// dateTimeFormat {string}  target input date time format, default: YYYY-MM-DDTHH:mm:ss+08:00
// dateTimeStr	  {string}  target input date time string
func ParseDateTimeStrToTimestamp(dateTimeFormat string, dateTimeStr string) (int64, error) {
	if IsEmptyString(dateTimeStr) {
		return 0, errors.New("invalid dateTime string")
	}
	if IsEmptyString(dateTimeFormat) {
		dateTimeFormat = DefaultDateTimeFormat
	}
	tt, err := time.ParseInLocation(dateTimeFormat, dateTimeStr, loc)
	if err != nil {
		return 0, err
	}
	return tt.Unix(), err
}

// make current date time timestamp
func MakeTimestamp() int64 {
	return time.Now().In(loc).UnixNano() / int64(time.Millisecond)
}
