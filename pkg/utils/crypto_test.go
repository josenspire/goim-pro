package utils

import (
	"fmt"
	"testing"
)

func TestNewSHA256(t *testing.T) {
	var password string = "1234567890"
	var salt string = "JAMES01"
	hashStr := NewSHA256(password, salt)

	// 3d9faa489bf4a363a67f6081b7106ae780496025efa52e45fd0428d15f5893ec
	fmt.Println(hashStr)
}
