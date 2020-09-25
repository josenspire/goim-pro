package utils

import (
	"fmt"
	"testing"
)

func TestGetFirstLetter(t *testing.T) {
	result := GetStringLetters("德玛西亚321QWE")
	result2 := GetStringLetters("1Q德玛西亚")
	result3 := GetStringLetters("m德玛123西11亚")
	fmt.Println(result)
	fmt.Println(result2)
	fmt.Println(result3)
}
