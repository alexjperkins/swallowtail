package formatter

import (
	"fmt"
	"strings"
	tradeengineproto "swallowtail/s.trade-engine/proto"
)

var (
	shortEmoji = ":chart_with_downwards_trend:"
	longEmoji  = ":chart_with_upwards_trend:"
)

// FormatTrade ...
func FormatTrade(header string, trade *tradeengineproto.Trade, parsedFrom string) string {
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

CURRENT_PRICE %s

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
		trade.ActorId,
		trade.ActorType.String(),
		trade.CurrentPrice,
		trade.Entry,
		trade.StopLoss,
	)

	// Append take profits if they exist.
	var footer strings.Builder
	for i, tp := range trade.TakeProfits {
		footer.WriteString(fmt.Sprintf("TP%v:        %v\n", i+1, tp))
	}

	// Append where we parsed the trade from.
	footer.WriteString(fmt.Sprintf("\nParsed From:\n%s", parsedFrom))

	return fmt.Sprintf("%s%s```%s%s```", base, warning, formattedContent, footer.String())
}
