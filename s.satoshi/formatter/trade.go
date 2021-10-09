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

// FormatTrade humanizes a trade in string format.
func FormatTrade(header string, trade *tradeengineproto.Trade, parsedFrom string) string {
	switch len(trade.Entries) {
	case 0:
		return formatTrade(header, trade, parsedFrom)
	case 1:
		return formatTrade(header, trade, parsedFrom)
	default:
		return formatDCATrade(header, trade, parsedFrom)
	}
}

func formatTrade(header string, trade *tradeengineproto.Trade, parsedFrom string) string {
	var sideEmoji string
	switch trade.TradeSide {
	case tradeengineproto.TRADE_SIDE_LONG, tradeengineproto.TRADE_SIDE_BUY:
		sideEmoji = longEmoji
	case tradeengineproto.TRADE_SIDE_SHORT, tradeengineproto.TRADE_SIDE_SELL:
		sideEmoji = shortEmoji
	}

	base := fmt.Sprintf("%s   `NEW TRADE ALERT: %s: %s%s`    :rocket:", sideEmoji, header, trade.Asset, trade.Pair)
	warning := `

:warning: Satoshi can not and **will** not be 100% accurate; please make sure the trade is sensible before placing :warning:
`

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

CURRENT_PRICE %v

ENTRY:        %v
STOP LOSS:    %v
`
	formattedContent := fmt.Sprintf(
		content,
		trade.TradeId,
		trade.Created.AsTime(),

		strings.ToUpper(trade.Asset),
		trade.Pair.String(),

		trade.TradeType.String(),
		trade.TradeSide.String(),
		trade.OrderType.String(),
		trade.HumanizedActorName,
		trade.ActorType.String(),
		trade.CurrentPrice,
		trade.Entries[0],
		trade.StopLoss,
	)

	// Append take profits if they exist.
	var footer strings.Builder
	for i, tp := range trade.TakeProfits {
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

func formatDCATrade(header string, trade *tradeengineproto.Trade, parsedFrom string) string {
	var sideEmoji string
	switch trade.TradeSide {
	case tradeengineproto.TRADE_SIDE_LONG, tradeengineproto.TRADE_SIDE_BUY:
		sideEmoji = longEmoji
	case tradeengineproto.TRADE_SIDE_SHORT, tradeengineproto.TRADE_SIDE_SELL:
		sideEmoji = shortEmoji
	}

	base := fmt.Sprintf("%s   `NEW DCA TRADE ALERT: %s: %s%s`    :lizard:", sideEmoji, header, trade.Asset, trade.Pair)
	warning := `

:warning: This is a DCA Order. Satoshi can not and **will** not be 100% accurate; please make sure the trade is sensible before placing :warning:
`
	sortedEntries := trade.Entries
	sort.Slice(sortedEntries, func(i, j int) bool {
		return sortedEntries[i] > sortedEntries[j]
	})

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

CURRENT_PRICE %v

UPPER:        %v
LOWER:        %v
STOP LOSS:    %v
`
	formattedContent := fmt.Sprintf(
		content,
		trade.TradeId,
		trade.Created.AsTime(),

		strings.ToUpper(trade.Asset),
		trade.Pair.String(),

		trade.TradeType.String(),
		trade.TradeSide.String(),
		trade.OrderType.String(),
		trade.HumanizedActorName,
		trade.ActorType.String(),
		trade.CurrentPrice,
		sortedEntries[0],
		sortedEntries[1],
		trade.StopLoss,
	)

	// Append take profits if they exist.
	var footer strings.Builder
	for i, tp := range trade.TakeProfits {
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
