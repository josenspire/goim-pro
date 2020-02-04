package utils

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewToken(t *testing.T) {
	Convey("Testing_create_new_token", t, func() {
		tokenStr := NewToken([]byte("13631210000"))
		fmt.Println(tokenStr)
		So(tokenStr, ShouldNotEqual, "")
	})
}
