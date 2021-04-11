package domain

import "time"

type Account struct {
	Username     string
	Password     string
	Transactions []*Asset
	assets       map[string]float64
	created      time.Time
	lastUpdated  time.Time
}
