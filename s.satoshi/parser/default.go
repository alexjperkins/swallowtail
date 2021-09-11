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

var (
	// TODO: we need proper gRPC service testing.
	fetchLatestPrice = getLatestPrice
)

func init() {
	register(discordproto.DiscordMoonModMessagesChannel, &DefaultParser{})
	register(discordproto.DiscordMoonSwingGroupChannel, &DefaultParser{})
	register(discordproto.DiscordSatoshiInternalCallsChannel, &DefaultParser{})
}

// The default parser to parse trades from external mods.
type DefaultParser struct{}

// Parse attempts to parse some content into a `tradeengineproto.Trade`. If it fails it returns a `FailedPrecondition` gerror
// that details why it was unable to.
func (d *DefaultParser) Parse(ctx context.Context, content string, m *discordgo.MessageCreate, actorType tradeengineproto.ACTOR_TYPE) (*tradeengineproto.Trade, error) {
	ticker := parseTicker(content)
	if ticker == "" {
		return nil, gerrors.FailedPrecondition("failed_to_parse_default.not_enough_information.missing_ticker", nil)
	}

	possibleValues, err := parseNumbersFromContent(content)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_parse_default", nil)
	}

	if len(possibleValues) < 2 {
		return nil, gerrors.FailedPrecondition("failed_to_parse_default.not_enough_information.values", map[string]string{
			"ticker": ticker,
		})
	}

	side, _ := parseSide(content)

	switch {
	case side == tradeengineproto.TRADE_SIDE_LONG:
		sort.Float64s(possibleValues)
	case side == tradeengineproto.TRADE_SIDE_SHORT:
		// Reverse Sort
		sort.Slice(possibleValues, func(i, j int) bool {
			return possibleValues[i] > possibleValues[j]
		})
	}

	currentPrice, err := fetchLatestPrice(ctx, ticker, tradeengineproto.TRADE_PAIR_USD.String())
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_parse_default", nil)
	}

	entry, stopLoss, takeProfits, err := validatePosition(ticker, "USDT", currentPrice, possibleValues)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_parse_with_default_parser.validate_position", nil)
	}

	var orderType, _ = parseOrderType(content, currentPrice, entry, side)

	protoTakeProfits := make([]float32, 0, len(takeProfits))
	for _, tp := range takeProfits {
		protoTakeProfits = append(protoTakeProfits, float32(tp))
	}

	actor := parseActor(m.Author.Username)

	return &tradeengineproto.Trade{
		ActorId:      actor,
		ActorType:    actorType,
		Asset:        strings.ToUpper(ticker),
		Pair:         tradeengineproto.TRADE_PAIR_USDT,
		TradeSide:    side,
		CurrentPrice: float32(currentPrice),
		Entry:        float32(entry),
		StopLoss:     float32(stopLoss),
		TakeProfits:  protoTakeProfits,
		OrderType:    orderType,
		TradeType:    tradeengineproto.TRADE_TYPE_FUTURES_PERPETUALS,
	}, nil
}
