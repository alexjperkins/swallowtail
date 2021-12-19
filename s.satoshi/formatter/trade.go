package formatter

import (
	"fmt"
	"sort"
	"strings"
	tradeengineproto "swallowtail/s.trade-engine/proto"
)

var (
	shortEmoji = ":chart_with_downwards_trend:"
	longEmoji  = ":chart_with_upwards_trend:"
)

// FormatTradeStrategy humanizes a trade in string format.
func FormatTradeStrategy(header string, trade *tradeengineproto.TradeStrategy, parsedFrom string) string {
	switch len(trade.Entries) {
	case 0:
		return formatTradeStrategy(header, trade, parsedFrom)
	case 1:
		return formatTradeStrategy(header, trade, parsedFrom)
	default:
		return formatDCATrade(header, trade, parsedFrom)
	}
}

func formatTradeStrategy(header string, tradeStrategy *tradeengineproto.TradeStrategy, parsedFrom string) string {
	var sideEmoji string
	switch tradeStrategy.TradeSide {
	case tradeengineproto.TRADE_SIDE_LONG, tradeengineproto.TRADE_SIDE_BUY:
		sideEmoji = longEmoji
	case tradeengineproto.TRADE_SIDE_SHORT, tradeengineproto.TRADE_SIDE_SELL:
		sideEmoji = shortEmoji
	}

	base := fmt.Sprintf("%s   `NEW TRADE ALERT: %s: %s%s`    :rocket:", sideEmoji, header, tradeStrategy.Asset, tradeStrategy.Pair)
	warning := `

:warning: Satoshi can not and **will** not be 100% accurate; please make sure the trade is sensible before placing :warning:
`

	var venues []string
	for _, v := range tradeStrategy.TradeableVenues {
		venues = append(venues, v.String())
	}

	content := `
TRADE ID:     %s 
TIMESTAMP:    %v

ASSET:        %v
PAIR:         %v
TRADE TYPE:   %s
TRADE SIDE:   %s
ORDER TYPE:   %s
MOD:          %s
MOD TYPE:     %s
VENUES:       %s

CURRENT_PRICE %v

ENTRY:        %v
STOP LOSS:    %v
`
	formattedContent := fmt.Sprintf(
		content,
		tradeStrategy.TradeStrategyId,
		tradeStrategy.Created.AsTime(),

		strings.ToUpper(tradeStrategy.Asset),
		tradeStrategy.Pair.String(),

		tradeStrategy.InstrumentType.String(),
		tradeStrategy.TradeSide.String(),
		tradeStrategy.ExecutionStrategy.String(),
		tradeStrategy.HumanizedActorName,
		tradeStrategy.ActorType.String(),
		strings.Join(venues, ", "),
		tradeStrategy.CurrentPrice,
		tradeStrategy.Entries[0],
		tradeStrategy.StopLoss,
	)

	// Append take profits if they exist.
	var footer strings.Builder
	for i, tp := range tradeStrategy.TakeProfits {
		footer.WriteString(fmt.Sprintf("TP%v:          %v\n", i+1, tp))
	}

	riskMessage := `
Please manage your risk accordingly. To **place** a trade react with one of the following emojis within **15 minutes**:

1%:  :one:
2%:  :two:
5%:  :five:
10%: :keycap_ten:

Always manually check the trade has been put on correctly on your account. Don't assume it will work 100% of the time whilst in **Beta**.
`
	// Append where we parsed the trade from.
	footer.WriteString(fmt.Sprintf("\nParsed From:\n%s", parsedFrom))

	return fmt.Sprintf("%s%s```%s%s```%s", base, warning, formattedContent, footer.String(), riskMessage)
}

func formatDCATrade(header string, tradeStrategy *tradeengineproto.TradeStrategy, parsedFrom string) string {
	var sideEmoji string
	switch tradeStrategy.TradeSide {
	case tradeengineproto.TRADE_SIDE_LONG, tradeengineproto.TRADE_SIDE_BUY:
		sideEmoji = longEmoji
	case tradeengineproto.TRADE_SIDE_SHORT, tradeengineproto.TRADE_SIDE_SELL:
		sideEmoji = shortEmoji
	}

	base := fmt.Sprintf("%s   `NEW DCA TRADE ALERT: %s: %s%s`    :lizard:", sideEmoji, header, tradeStrategy.Asset, tradeStrategy.Pair)
	warning := `

:warning: This is a DCA Order. Satoshi can not and **will** not be 100% accurate; please make sure the trade is sensible before placing :warning:
`
	sortedEntries := tradeStrategy.Entries
	sort.Slice(sortedEntries, func(i, j int) bool {
		return sortedEntries[i] > sortedEntries[j]
	})

	var venues []string
	for _, v := range tradeStrategy.TradeableVenues {
		venues = append(venues, v.String())
	}

	content := `
TRADE ID:     %s 
TIMESTAMP:    %v

ASSET:        %v
PAIR:         %v
TRADE TYPE:   %s
TRADE SIDE:   %s
ORDER TYPE:   %s
MOD:          %s
MOD TYPE:     %s
VENUEs:       %s

CURRENT_PRICE %v

UPPER:        %v
LOWER:        %v
STOP LOSS:    %v
`
	formattedContent := fmt.Sprintf(
		content,
		tradeStrategy.TradeStrategyId,
		tradeStrategy.Created.AsTime(),

		strings.ToUpper(tradeStrategy.Asset),
		tradeStrategy.Pair.String(),

		tradeStrategy.InstrumentType.String(),
		tradeStrategy.TradeSide.String(),
		tradeStrategy.ExecutionStrategy.String(),
		tradeStrategy.HumanizedActorName,
		tradeStrategy.ActorType.String(),
		strings.Join(venues, ", "),
		tradeStrategy.CurrentPrice,
		sortedEntries[0],
		sortedEntries[1],
		tradeStrategy.StopLoss,
	)

	// Append take profits if they exist.
	var footer strings.Builder
	for i, tp := range tradeStrategy.TakeProfits {
		footer.WriteString(fmt.Sprintf("TP%v:          %v\n", i+1, tp))
	}

	riskMessage := `
Please manage your risk accordingly. To **place** a trade react with one of the following emojis within **15 minutes**:

1%:  :one:
2%:  :two:
5%:  :five:
10%: :keycap_ten:

Always manually check the trade has been put on correctly on your account. Don't assume it will work 100% of the time whilst in **Beta**.
`
	// Append where we parsed the trade from.
	footer.WriteString(fmt.Sprintf("\nParsed From:\n%s", parsedFrom))

	return fmt.Sprintf("%s%s```%s%s```%s", base, warning, formattedContent, footer.String(), riskMessage)
}
