package handler

import (
	"math/rand"
	"time"
)

func jitter(min, max int) time.Duration {
	rand.Seed(time.Now().UTC().UnixNano())
	return time.Duration(min+rand.Intn(max)) * time.Second
}

func abs(a float64) float64 {
	if a < 0 {
		return -a
	}
	return a
}
