package utils

import (
	"fmt"
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
