package exchange

import (
	"context"
	"strings"
	"swallowtail/libraries/gerrors"
	accountproto "swallowtail/s.account/proto"
	"swallowtail/s.trade-engine/domain"
	"time"
)

// ExecuteTrade ...
func ExecuteTrade(ctx context.Context, trade *domain.Trade, requestTimestamp time.Time) error {
	switch strings.ToUpper(trade.Exchange) {
	case accountproto.ExchangeType_FTX.String():
		return gerrors.Unimplemented("failed_to_execute_trade.exchange_not_supported", map[string]string{
			"exchange": trade.Exchange,
		})
	case accountproto.ExchangeType_BINANCE.String():
		return executeBinanceTrade(ctx, trade, requestTimestamp)
	default:
		return gerrors.Unimplemented("failed_to_execute_trade.exchange_not_supported", map[string]string{
			"exchange": trade.Exchange,
		})
	}
}
