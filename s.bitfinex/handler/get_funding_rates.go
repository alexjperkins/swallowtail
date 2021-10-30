package handler

import (
	"context"

	"swallowtail/libraries/gerrors"
	"swallowtail/s.bitfinex/client"
	"swallowtail/s.bitfinex/dto"
	"swallowtail/s.bitfinex/marshaling"
	bitfinexproto "swallowtail/s.bitfinex/proto"
)

// GetBitfinexFundingRates fetches the funding rates from Bitfinex by `symbol` & `limit`.
func (s *BitfinexService) GetBitfinexFundingRates(
	ctx context.Context, in *bitfinexproto.GetBitfinexFundingRatesRequest,
) (*bitfinexproto.GetBitfinexFundingRatesResponse, error) {
	// Validation.
	switch {
	case in.Symbol == "":
		return nil, gerrors.BadParam("missing_param.symbol", nil)
	}

	errParams := map[string]string{
		"symbol": in.Symbol,
	}

	// Get funding rate.
	rsp, err := client.GetFundingRates(ctx, &dto.GetFundingRatesRequest{
		Symbol: in.Symbol,
	})
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_get_funding_rates", errParams)
	}

	// Marshal to proto.
	return marshaling.GetFundingRatesDTOToProto(rsp), nil
}
