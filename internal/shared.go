package internal

import (
	"math/rand"
	"time"
)

func RandInt(max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max)
}
