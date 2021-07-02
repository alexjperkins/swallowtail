package domain

import "time"

// Touch represents a new touch for a message & will be used for idempotency reasons.
type Touch struct {
	IdempotencyKey string    `db:"idempotency_key"`
	Updated        time.Time `db:"timestamp"`
	SenderID       string    `db:"sender_id"`
}
