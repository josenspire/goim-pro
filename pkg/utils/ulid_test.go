package utils

import (
	"fmt"
	"testing"
)

func TestNewULID(t *testing.T) {
	ulid := NewULID()
	fmt.Println(ulid)
}
