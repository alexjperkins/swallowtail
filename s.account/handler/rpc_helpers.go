package handler

import (
	"context"
	"strings"
	"swallowtail/libraries/gerrors"
	"swallowtail/s.account/domain"
)

func validateExchangeCredentials(ctx context.Context, exchange *domain.Exchange) (bool, string, error) {
	errParams := map[string]string{
		"exchange_type": exchange.Exchange,
	}
	switch strings.ToLower(exchange.Exchange) {
	case "binance":
		return validateBinanceExchangeCredentials(ctx, exchange)
	case "ftx":
		return validateFTXExchangeCredentials(ctx, exchange)
	default:
		return false, "", gerrors.FailedPrecondition("invalid_exchangge", errParams)
	}
}

func validateBinanceExchangeCredentials(ctx context.Context, exchange *domain.Exchange) (bool, string, error) {
	return false, "", nil
}

func validateFTXExchangeCredentials(ctx context.Context, exchange *domain.Exchange) (bool, string, error) {
	return false, "", nil
}
