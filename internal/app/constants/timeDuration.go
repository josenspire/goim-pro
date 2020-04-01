package constants

import "time"

const (
	DefaultExpiresTime = time.Minute * time.Duration(30) // 30 min

	SecondOf60 = time.Second * time.Duration(60) // 60s

	MinuteOf15 = time.Minute * time.Duration(15) // 15 min

	ThreeDays = time.Hour * time.Duration(24*3) // 3 days
)
