package domain

import "time"

type Account struct {
	UserID               string    `db:"user_id"`
	Username             string    `db:"username"`
	Password             string    `db:"password"`
	Email                string    `db:"email"`
	PhoneNumber          string    `db:"phone_number"`
	PrimaryVenueAccount  string    `db:"primary_venue_account"`
	Created              time.Time `db:"created"`
	Updated              time.Time `db:"updated"`
	LastPaymentTimestamp time.Time `db:"last_payment_timestamp"`
	HighPriorityPager    string    `db:"high_priority_pager"`
	LowPriorityPager     string    `db:"low_priority_pager"`
	IsAdmin              bool      `db:"is_admin"`
	IsFuturesMember      bool      `db:"is_futures_member"`
	DefaultDCAStrategy   string    `db:"default_dca_strategy"`
}
