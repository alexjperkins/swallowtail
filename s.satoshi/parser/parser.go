package parser

import (
	tradeengineoroto "swallowtail/s.trade-engine/proto"

	"github.com/bwmarrin/discordgo"
)

// TradeParser ...
type TradeParser interface {
	Parse(content string, m *discordgo.MessageCreate) (*tradeengineoroto.Trade, bool)
}

// Parse ...
func Parse(identifier, content string, m *discordgo.MessageCreate) (*tradeengineoroto.Trade, bool) {
	parser, ok := getParserByIdentifier(identifier)
	if !ok {
		return nil, false
	}

	return parser.Parse(content, m)
}
