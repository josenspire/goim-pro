package utils

import "time"

const (
	defaultDateTimeFormat = "2019-01-08 13:50:30"
)

func ParseTimestampToDateTime(timestamp int64) time.Time {
	if timestamp == 0 {
		return time.Unix(0, 0)
	}
	return time.Unix(timestamp, 0)
}

func ParseTimestampToDateTimeStr(timestamp int64, format string) string {
	if timestamp == 0 {
		return ""
	}
	if format == "" {
		format = defaultDateTimeFormat
	}
	return time.Unix(timestamp, 0).Format(format)
}
