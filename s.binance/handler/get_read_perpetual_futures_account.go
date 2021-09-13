package handler

import (
	"context"
	"strings"

	"swallowtail/libraries/gerrors"
	"swallowtail/s.binance/client"
	"swallowtail/s.binance/marshaling"
	binanceproto "swallowtail/s.binance/proto"
)

// ReadPerpetualFuturesAccount ...
func (s *BinanceService) ReadPerpetualFuturesAccount(
	ctx context.Context, in *binanceproto.ReadPerpetualFuturesAccountRequest,
) (*binanceproto.ReadPerpetualFuturesAccountResponse, error) {
	switch {
	case in.ActorId == "":
		return nil, gerrors.BadParam("missing_param.actor_id", nil)
	case !isValidActor(in.ActorId):
		return nil, gerrors.Unauthenticated("failed_to_read_perpetual_futures_account.unauthorized", nil)
	}

	if err := isValidCredentials(in.Credentials, false); err != nil {
		return nil, gerrors.Augment(err, "failed_to_read_perpetual_futures_account.credentials", nil)
	}

	errParams := map[string]string{
		"actor_id": in.ActorId,
	}

	rsp, err := client.ReadPerpetualFuturesAccount(ctx, nil, &client.Credentials{
		APIKey:    in.GetCredentials().ApiKey,
		SecretKey: in.GetCredentials().SecretKey,
	})
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_read_perpetual_futures_account", errParams)
	}

	var usdtBalance *client.PerpetualFuturesAccountBalance
	for _, balance := range *rsp {
		if strings.ToLower(balance.Asset) != "usdt" {
			continue
		}

		usdtBalance = balance
	}

	if usdtBalance == nil {
		return nil, gerrors.NotFound("failed_to_read_perpetual_futures_account.usdt_balance_not_found_futures_account", errParams)
	}

	protoRsp, err := marshaling.PerpetualFuturesAccountBalanceDTOToProto(usdtBalance)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_read_perpetual_futures_account.marshal_balance_to_proto", errParams)
	}

	return protoRsp, err
}
