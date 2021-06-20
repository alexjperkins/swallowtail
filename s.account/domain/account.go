package domain

type Account struct {
	AccountID         string `db:"account_id"`
	Username          string `db:"username"`
	Password          string `db:"password"`
	Email             string `db:"email"`
	DiscordID         string `db:"discord_id"`
	PhoneNumber       string `db:"phone_number"`
	HighPriorityPager string `db:"high_priority_pager"`
	LowPriorityPager  string `db:"low_priority_pager"`
}
