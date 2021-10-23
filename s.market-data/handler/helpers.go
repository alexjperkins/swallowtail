package handler

import (
	"math/rand"
	"time"
)

func jitter(min, max int) time.Duration {
	rand.Seed(time.Now().UTC().UnixNano())
	return time.Duration(min+rand.Intn(max)) * time.Second
}
