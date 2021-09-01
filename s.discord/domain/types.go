package domain

import "time"

// Touch represents a new touch for a message & will be used for idempotency reasons.
type Touch struct {
	TouchID        string    `db:"touch_id"`
	IdempotencyKey string    `db:"idempotency_key"`
	Updated        time.Time `db:"updated"`
	SenderID       string    `db:"sender_id"`
}

// Role ...
type Role struct {
	ID   string
	Name string
}
