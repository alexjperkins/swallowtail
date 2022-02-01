package domain

import "time"

// Account defines the account object stored in the domain.
type Account struct {
	UserID               string    `db:"user_id"`
	Username             string    `db:"username"`
	Email                string    `db:"email"`
	PhoneNumber          string    `db:"phone_number"`
	HighPriorityPager    string    `db:"high_priority_pager"`
	LowPriorityPager     string    `db:"low_priority_pager"`
	Created              time.Time `db:"created"`
	Updated              time.Time `db:"updated"`
	LastPaymentTimestamp time.Time `db:"last_payment_timestamp"`
	PrimaryVenue         string    `db:"primary_venue"`
	IsAdmin              bool      `db:"is_admin"`
	IsFuturesMember      bool      `db:"is_futures_member"`
	DefaultDCAStrategy   string    `db:"default_dca_strategy"`
}
