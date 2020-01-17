package utils

import (
	"github.com/oklog/ulid"
	"math/rand"
	"time"
)

func NewULID() (newULID string) {
	t := time.Unix(time.Now().Unix(), 0)
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	newULID = ulid.MustNew(ulid.Timestamp(t), entropy).String()
	logger.Infof("[ulid]: %s", newULID)
	return
}
