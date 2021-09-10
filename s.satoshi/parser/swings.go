package parser

import (
	tradeengineproto "swallowtail/s.trade-engine/proto"

	"github.com/bwmarrin/discordgo"
)

const (
	swingsParserID = "swings-trade-parser"
)

func init() {
}

type SwingsParser struct{}

func (s *SwingsParser) Parse(content string, m *discordgo.MessageCreate) (*tradeengineproto.Trade, bool) {
	return nil, false
}
