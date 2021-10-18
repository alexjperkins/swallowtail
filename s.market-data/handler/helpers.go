package handler

import (
	"math/rand"
	"time"
)

func jitter(min, max int) time.Duration {
	return time.Duration(min+rand.Intn(max)) * time.Second
}
