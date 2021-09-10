package parser

import (
	"context"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"

	"swallowtail/libraries/gerrors"
	discordproto "swallowtail/s.discord/proto"
	tradeengineproto "swallowtail/s.trade-engine/proto"
)

func init() {
	register(discordproto.DiscordMoonModMessagesChannel, &WWGParser{})
}

type WWGParser struct{}

func (w *WWGParser) Parse(ctx context.Context, content string, m *discordgo.MessageCreate) (*tradeengineproto.Trade, error) {
	ticker := parseTicker(content)
	if ticker == "" {
		return nil, gerrors.FailedPrecondition("failed_to_parse_wwg.not_enough_information.missing_ticker", nil)
	}

	possibleValues, err := parseNumbersFromContent(content)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_parse_wwg", nil)
	}

	if len(possibleValues) < 2 {
		return nil, gerrors.FailedPrecondition("failed_to_parse_wwg.not_enough_information.values", map[string]string{
			"ticker": ticker,
		})
	}

	side, _ := parseSide(content)

	switch {
	case side == tradeengineproto.TRADE_SIDE_BUY:
		sort.Float64s(possibleValues)
	case side == tradeengineproto.TRADE_SIDE_SELL:
		// Reverse Sort
		sort.Slice(possibleValues, func(i, j int) bool {
			return possibleValues[i] > possibleValues[j]
		})
	}

	truth, err := fetchLatestPrice(ctx, ticker, tradeengineproto.TRADE_PAIR_USD.String())
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_parse_wwg", nil)
	}

	entry, stopLoss, takeProfits, err := validatePosition(ticker, "USDT", truth, possibleValues)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_parse_wwg.validate_position", nil)
	}

	var orderType, _ = parseOrderType(content, truth, entry, side)

	protoTakeProfits := make([]float32, 0, len(takeProfits))
	for _, tp := range takeProfits {
		protoTakeProfits = append(protoTakeProfits, float32(tp))
	}

	return &tradeengineproto.Trade{
		ActorId:     strings.ToUpper(m.Author.Username),
		ActorType:   tradeengineproto.ACTOR_TYPE_EXTERNAL.String(),
		Asset:       strings.ToUpper(ticker),
		Pair:        tradeengineproto.TRADE_PAIR_USDT,
		TradeSide:   side,
		Entry:       float32(entry),
		StopLoss:    float32(stopLoss),
		TakeProfits: protoTakeProfits,
		OrderType:   orderType,
		TradeType:   tradeengineproto.TRADE_TYPE_FUTURES_PERPETUALS,
	}, nil
}
