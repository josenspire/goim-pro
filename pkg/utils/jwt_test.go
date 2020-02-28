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
		So(tokenStr, ShouldNotEqual, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJNVE0yTXpFeU1UQXdNREE9IiwiZXhwIjoxNTgwODY5NTMyLCJpYXQiOjE1ODA4NjU5MzIsImlzcyI6InNhbHR5X2ltIn0")
	})
}

func TestTokenVerify(t *testing.T) {
	Convey("Testing_verify_token", t, func() {
		Convey("Token_Should_Be_Valid", func() {
			tokenStr := NewToken([]byte("13631210000"))
			fmt.Println(tokenStr)

			isValid, payload, err := TokenVerify(tokenStr)
			if err != nil {
				t.FailNow()
			}
			So(isValid, ShouldBeTrue)
			So(string(payload), ShouldEqual, "13631210000")
		})

		Convey("Token_Should_Be_InValid", func() {
			var tokenStr string = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJNVE0yTXpFeU1UQXdNREE9IiwiZXhwIjoxNTgwODY3MTU4LCJpYXQiOjE1ODA4NjcxNTgsImlzcyI6InNhbHR5X2ltIn0.F_a6DS5QHo5JfbO1pMo_WYLuT4X_sUF-V9on2b2L7Lc`
			fmt.Println(tokenStr)

			isValid, _, err := TokenVerify(tokenStr)

			So(isValid, ShouldBeFalse)
			So(err.Error(), ShouldNotBeEmpty)
		})
	})
}