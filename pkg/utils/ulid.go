package utils

import (
	"github.com/oklog/ulid"
	"math/rand"
	"time"
)

func NewULID(seed int64) (newULID string) {
	sec := time.Now().Unix()
	if seed != 0 {
		sec += seed
	}
	t := time.Unix(sec, 0)
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	newULID = ulid.MustNew(ulid.Timestamp(t), entropy).String()
	logger.Infof("[ulid]: %s", newULID)
	return
}
