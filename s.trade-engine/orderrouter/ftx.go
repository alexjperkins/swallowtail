package orderrouter

import (
	"context"

	"swallowtail/libraries/gerrors"
	ftxproto "swallowtail/s.ftx/proto"
	"swallowtail/s.trade-engine/domain"
	tradeengineproto "swallowtail/s.trade-engine/proto"
)

func executeFTXFuturesTrade(
	ctx context.Context,
	trade *domain.Trade,
	participant *tradeengineproto.AddParticipantToTradeRequest,
	credentials *ftxproto.FTXCredentials,
) (*FuturesTradeResponse, error) {
	return nil, gerrors.Unimplemented("faield_to_execute_ftx_futures_trade.unimplemented", nil)
}
