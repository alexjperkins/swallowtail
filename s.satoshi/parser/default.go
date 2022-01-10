package parser

import (
	"context"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"

	"swallowtail/libraries/gerrors"
	tradeengineproto "swallowtail/s.trade-engine/proto"
)

var (
	// TODO: we need proper gRPC service testing.
	fetchLatestPrice = getLatestPrice
)

// The default parser to parse trades from external mods.
type DefaultParser struct{}

// Parse attempts to parse some content into a `tradeengineproto.Trade`. If it fails it returns a `FailedPrecondition` gerror
// that details why it was unable to.
func (d *DefaultParser) Parse(ctx context.Context, content string, m *discordgo.MessageCreate, actorType tradeengineproto.ACTOR_TYPE) (*tradeengineproto.TradeStrategy, error) {
	// Parse instrument types.
	instrumentType := parseInstrumentType(content)

	// Parse venues.
	ticker, venues := parseTickerAndVenues(content, instrumentType)
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

	errParams := map[string]string{
		"ticker": ticker,
	}

	currentPrice, err := fetchLatestPrice(ctx, ticker)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_parse_default", errParams)
	}

	entries, stopLoss, takeProfits, err := validatePosition(currentPrice, possibleValues, false)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_parse_with_default_parser.validate_position", nil)
	}

	if len(entries) != 1 {
		errParams["entries"] = entriesAsString(entries)
		return nil, gerrors.FailedPrecondition("failed_to_parse_with_default_parser.multiple_entries", errParams)
	}

	executionStrategy, _ := parseExecutionStrategy(content, currentPrice, entries, side)

	protoEntries := make([]float32, 0, len(entries))
	for _, entry := range entries {
		protoEntries = append(protoEntries, float32(entry))
	}

	protoTakeProfits := make([]float32, 0, len(takeProfits))
	for _, tp := range takeProfits {
		protoTakeProfits = append(protoTakeProfits, float32(tp))
	}

	actor := parseActor(m.Author.Username)

	return &tradeengineproto.TradeStrategy{
		ActorId:            m.Author.ID,
		HumanizedActorName: actor,
		ActorType:          actorType,
		ExecutionStrategy:  executionStrategy,
		InstrumentType:     instrumentType,
		TradeSide:          side,
		Asset:              strings.ToUpper(ticker),
		Pair:               tradeengineproto.TRADE_PAIR_USDT,
		Entries:            protoEntries,
		StopLoss:           float32(stopLoss),
		TakeProfits:        protoTakeProfits,
		CurrentPrice:       float32(currentPrice),
		TradeableVenues:    venues,
	}, nil
}
