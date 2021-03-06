package utils

import (
	"math/rand"
	"time"
)

const StrSource string = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(time.Now().Unix()))
}

func GenerateRandString(length int) string {
	bytes := []byte(StrSource)
	var result []byte
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func IsEmptyStrings(strArgs ...string) bool {
	var result = true
	for _, str := range strArgs {
		if str != "" {
			result = false
			break
		}
	}
	return result
}

func IsContainEmptyString(strArgs ...string) bool {
	var result = false
	for _, str := range strArgs {
		if str == "" {
			result = true
			break
		}
	}
	return result
}
