package internal

import (
	"math/rand"
	"time"
)

func getShiftValue(maxStep int) int {
	rand.Seed(time.Now().UnixNano())

	return rand.Intn(maxStep)
}
