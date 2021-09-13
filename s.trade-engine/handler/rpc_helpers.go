package handler

import (
	"context"

	"swallowtail/libraries/gerrors"
	accountproto "swallowtail/s.account/proto"
	binanceproto "swallowtail/s.binance/proto"
)

func readPrimaryExchangeCredentials(ctx context.Context, userID string) (*accountproto.Exchange, error) {
	rsp, err := (&accountproto.ReadPrimaryExchangeByUserIDRequest{
		UserId:  userID,
		ActorId: "system:tradeengine",
	}).Send(ctx).Response()
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_read_primary_account_credentials", nil)
	}

	return rsp.PrimaryExchange, nil
}

func readBinanceFuturesAccountBalance(ctx context.Context, binanceCredentials *binanceproto.Credentials) (*binanceproto.ReadPerpetualFuturesAccountResponse, error) {
	rsp, err := (&binanceproto.ReadPerpetualFuturesAccountRequest{
		Credentials: binanceCredentials,
		ActorId:     binanceproto.BinanceAccountActorTradeEngineSystem,
	}).Send(ctx).Response()
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_read_binance_account_futures_balance", nil)
	}

	return rsp, nil
}
