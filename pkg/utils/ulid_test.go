package utils

import (
	"fmt"
	"testing"
)

func TestNewULID(t *testing.T) {
	ulid1 := NewULID(0)
	ulid2 := NewULID(1)
	fmt.Println(ulid1, ulid2)
}

func TestNewUserULID(t *testing.T) {
	ulid1 := NewUserULID(0)
	ulid2 := NewUserULID(1)
	fmt.Println(ulid1, ulid2)
}
