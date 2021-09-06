package parser

import (
	tradeengineoroto "swallowtail/s.trade-engine/proto"

	"github.com/bwmarrin/discordgo"
)

const (
	wwgParserID = "wwg-trade-parser"
)

func init() {
	register(wwgParserID, &WWGParser{})
}

type WWGParser struct{}

func (w *WWGParser) Parse(content string, m *discordgo.MessageCreate) (*tradeengineoroto.Trade, bool) {
	return nil, nil
}
