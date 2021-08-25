package handler

import (
	"context"
	"swallowtail/libraries/gerrors"
	accountproto "swallowtail/s.account/proto"
	binanceproto "swallowtail/s.binance/proto"
	"time"

	"github.com/monzo/slog"
)

func validateExchangeCredentials(ctx context.Context, userID string, exchange *accountproto.Exchange) (bool, string, error) {
	errParams := map[string]string{
		"exchange_type": exchange.ExchangeType.String(),
	}

	switch exchange.ExchangeType.String() {
	case accountproto.ExchangeType_BINANCE.String():
		return validateBinanceExchangeCredentials(ctx, userID, exchange)
	case accountproto.ExchangeType_FTX.String():
		return validateFTXExchangeCredentials(ctx, userID, exchange)
	default:
		return false, "", gerrors.FailedPrecondition("failed_to_validate_credentials.invalid_exchange", errParams)
	}
}

func validateBinanceExchangeCredentials(ctx context.Context, userID string, exchange *accountproto.Exchange) (bool, string, error) {
	slog.Warn(ctx, "Creds: %v", exchange)

	rsp, err := (&binanceproto.VerifyCredentialsRequest{
		UserId: userID,
		Credentials: &binanceproto.Credentials{
			ApiKey:    exchange.ApiKey,
			SecretKey: exchange.SecretKey,
		},
	}).SendWithTimeout(ctx, 30*time.Second).Response()
	if err != nil {
		return false, "", gerrors.Augment(err, "failed_to_validate_binance_credentials", nil)
	}

	return rsp.Success, rsp.Reason, nil
}

func validateFTXExchangeCredentials(ctx context.Context, userID string, exchange *accountproto.Exchange) (bool, string, error) {
	return false, "", nil
}
