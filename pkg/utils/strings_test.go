package utils

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGenerateRandomNum(t *testing.T) {
	result := GenerateRandomNum(6)
	result1 := GenerateRandomNum(6)
	result2 := GenerateRandomNum(6)
	result3 := GenerateRandomNum(6)
	result4 := GenerateRandomNum(6)
	fmt.Println(result, result1, result2, result3, result4)
}

func TestStrArray2String(t *testing.T) {
	strArr := []string{
		"136", "134", "137",
	}

	result := StrArray2String(strArr, ";")
	Convey("Test_Array2String", t, func() {
		Convey("should_split_with_;_and_return_actual_result", func() {
			So(result, ShouldEqual, "136;134;137")
		})
	})
}

func TestStrArrayDeduplication(t *testing.T) {
	var strArr = []string{
		"a1", "a2", "a3", "a1", "a2",
	}

	uniqStrArr := StrArrayDeduplication(strArr)
	Convey("should_deduplication_successful_then_return", t, func() {
		So(len(uniqStrArr), ShouldEqual, 3)
		So(uniqStrArr[0], ShouldEqual, "a1")
		So(uniqStrArr[1], ShouldEqual, "a2")
		So(uniqStrArr[2], ShouldEqual, "a3")
	})
}
