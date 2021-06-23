package accountproto

type PagerType string

const (
	AccountPagerTypeDiscord PagerType = "discord"
	AccountPagerTypeEmail             = "email"
	AccountPagerTypeSMS               = "sms"
	AccountPagerTypePhone             = "phone"
	AccountPagerTypeUnknown           = "unknown"
)
