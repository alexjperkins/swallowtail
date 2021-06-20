package domain

import "time"

type Account struct {
	AccountID            string    `db:"account_id"`
	Username             string    `db:"username"`
	Password             string    `db:"password"`
	Email                string    `db:"email"`
	DiscordID            string    `db:"discord_id"`
	PhoneNumber          string    `db:"phone_number"`
	Created              time.Time `db:"created"`
	Updated              time.Time `db:"updated"`
	LastPaymentTimestamp time.Time `db:"last_payment_timestamp"`
	HighPriorityPager    string    `db:"high_priority_pager"`
	LowPriorityPager     string    `db:"low_priority_pager"`
	IsAdmin              bool      `db:"is_admin"`
	IsFuturesMember      bool      `db:"is_futures_member"`
}
