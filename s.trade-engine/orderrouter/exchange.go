package orderrouter

import (
	"context"
	"time"

	"swallowtail/libraries/gerrors"
	accountproto "swallowtail/s.account/proto"
	binanceproto "swallowtail/s.binance/proto"
	ftxproto "swallowtail/s.ftx/proto"
	"swallowtail/s.trade-engine/domain"
	tradeengineproto "swallowtail/s.trade-engine/proto"
)

// FuturesTradeResponse ...
type FuturesTradeResponse struct {
	ExchangeTradeID        string
	NotionalSize           float64
	ExecutionTimestamp     time.Time
	NumberOfExecutedOrders int
	ExecutionAlgoStrategy  string
}

// ExecuteFuturesTradeForParticipant ...
func ExecuteFuturesTradeForParticipant(
	ctx context.Context,
	trade *domain.Trade,
	participant *tradeengineproto.AddParticipantToTradeRequest,
	exchange *accountproto.Exchange,
) (*FuturesTradeResponse, error) {
	switch {
	case exchange.ExchangeType == accountproto.ExchangeType_BINANCE:
		return executeBinanceFuturesTrade(ctx, trade, participant, &binanceproto.Credentials{
			ApiKey:    exchange.ApiKey,
			SecretKey: exchange.SecretKey,
		})
	case exchange.ExchangeType == accountproto.ExchangeType_FTX:
		return executeFTXFuturesTrade(ctx, trade, participant, &ftxproto.FTXCredentials{
			ApiKey:     exchange.ApiKey,
			SecretKey:  exchange.SecretKey,
			Subaccount: exchange.SubAccount,
		})
	default:
		return nil, gerrors.Unimplemented("cannot_execute_trade.unimplemented_exchange", map[string]string{
			"exchange": exchange.ExchangeType.String(),
		})

	}
}
