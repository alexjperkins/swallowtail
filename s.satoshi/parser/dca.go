package parser

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"

	"swallowtail/libraries/gerrors"
	tradeengineproto "swallowtail/s.trade-engine/proto"
)

var (
	possibleStopLossMarks = []string{
		"sl",
		"stop",
		"stoploss",
	}
	possibleTakeProfitMarks = []string{
		"profit",
		"tp",
	}
)

type DCAParser struct{}

func (d *DCAParser) Parse(ctx context.Context, content string, m *discordgo.MessageCreate, actorType tradeengineproto.ACTOR_TYPE) (*tradeengineproto.Trade, error) {
	ticker := parseTicker(content)
	if ticker == "" {
		return nil, gerrors.FailedPrecondition("failed_to_parse_dca.missing_ticker", nil)
	}

	var (
		stopLossMark   string
		takeProfitMark string
	)
	for _, slm := range possibleStopLossMarks {
		if strings.Contains(content, slm) {
			stopLossMark = slm
			break
		}
	}
	for _, tpm := range possibleTakeProfitMarks {
		if strings.Contains(content, tpm) {
			takeProfitMark = tpm
			break
		}
	}

	// Clean up content if we have DCA orders; mods like to hypenate.
	var (
		entriesContent    = content
		stopLossContent   string
		takeProfitContent string
	)
	if stopLossContent != "" {
		stopLossSplits := strings.SplitAfter(content, stopLossMark)
		entriesContent = strings.ReplaceAll(stopLossSplits[0], "-", " ")

		switch {
		case takeProfitContent != "":
			takeProfitSplits := strings.Split(stopLossSplits[1], takeProfitMark)
			stopLossContent, takeProfitContent = takeProfitSplits[0], strings.ReplaceAll(takeProfitSplits[1], "-", "")
		default:
			stopLossContent = stopLossSplits[1]
		}
	}

	currentPrice, err := fetchLatestPrice(ctx, ticker)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_parse_dca", nil)
	}

	// Validate this is a DCA order; we do so by checking if we have `dca` in the content or we have at least
	// two entries in the parsed entry content.
	switch {
	case strings.Contains(entriesContent, "dca"):
	default:
		entries, err := parseNumbersFromContent(entriesContent)
		if err != nil {
			return nil, gerrors.Augment(err, "failed_to_parse_dca.failed_to_parse_entries_from_content", nil)
		}

		if len(entries) < 2 {
			return nil, gerrors.Augment(err, "failed_to_parse_dca.not_enough_entries", map[string]string{
				"entries": entriesAsString(entries),
			})
		}
	}

	possibleValues, err := parseNumbersFromContent(fmt.Sprintf("%s %s %s", entriesContent, stopLossContent, takeProfitContent))
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_parse_dca.failed_to_parse_values_from_content", nil)
	}

	side, _ := parseSide(content)

	switch {
	case side == tradeengineproto.TRADE_SIDE_LONG:
		sort.Float64s(possibleValues)
	case side == tradeengineproto.TRADE_SIDE_SHORT:
		// Reverse Sort if we are shorting.
		sort.Slice(possibleValues, func(i, j int) bool {
			return possibleValues[i] > possibleValues[j]
		})
	}

	entries, stopLoss, takeProfits, err := validatePosition(currentPrice, possibleValues, true)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_parse_with_dca_parser.validate_position", nil)
	}

	if len(entries) != 2 {
		return nil, gerrors.FailedPrecondition("failed_to_parse_with_dca_parser.not_enough_entries", map[string]string{
			"ticker":  ticker,
			"entries": entriesAsString(entries),
		})
	}

	orderType, _ := parseOrderType(content, currentPrice, entries, side)

	protoEntries := make([]float32, 0, len(entries))
	for _, entry := range entries {
		protoEntries = append(protoEntries, float32(entry))
	}

	protoTakeProfits := make([]float32, 0, len(takeProfits))
	for _, tp := range takeProfits {
		protoTakeProfits = append(protoTakeProfits, float32(tp))
	}

	actor := parseActor(m.Author.Username)

	return &tradeengineproto.Trade{
		ActorId:            m.Author.ID,
		HumanizedActorName: actor,
		ActorType:          actorType,
		OrderType:          orderType,
		TradeType:          tradeengineproto.TRADE_TYPE_FUTURES_PERPETUALS,
		TradeSide:          side,
		Asset:              strings.ToUpper(ticker),
		Pair:               tradeengineproto.TRADE_PAIR_USDT,
		Entries:            protoEntries,
		StopLoss:           float32(stopLoss),
		TakeProfits:        protoTakeProfits,
		CurrentPrice:       float32(currentPrice),
	}, nil
}
