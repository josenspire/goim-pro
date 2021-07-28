package wuid

import (
	"fmt"
	"testing"
)

func TestNewWUID(t *testing.T) {
	wuid1 := NewWUID()
	fmt.Println(wuid1)
	wuid2 := NewWUID()
	fmt.Println(wuid2)
}
