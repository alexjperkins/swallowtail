package emojis

type SatoshiRiskEmoji string

func (s SatoshiRiskEmoji) AsRiskPercentage() int {
	switch s {
	case EmojiUnicodeTen:
		return 10
	case EmojiUnicodeFive:
		return 5
	case EmojiUnicodeTwo:
		return 2
	case EmojiUnicodeOne:
		return 1
	default:
		return 0
	}
}

func (s SatoshiRiskEmoji) String() string {
	switch s {
	case EmojiUnicodeTen:
		return "10%"
	case EmojiUnicodeFive:
		return "5%"
	case EmojiUnicodeTwo:
		return "2%"
	case EmojiUnicodeOne:
		return "1%"
	default:
		return "unknown"
	}
}

const (
	// Satoshi specific unicode emojis
	EmojiUnicodeTen  SatoshiRiskEmoji = "🔟"
	EmojiUnicodeFive SatoshiRiskEmoji = "5️⃣"
	EmojiUnicodeTwo  SatoshiRiskEmoji = "2️⃣"
	EmojiUnicodeOne  SatoshiRiskEmoji = "1️⃣"
)
