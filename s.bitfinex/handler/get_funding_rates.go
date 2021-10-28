package handler

import (
	"context"
	"strconv"

	"swallowtail/libraries/gerrors"
	"swallowtail/s.bitfinex/client"
	"swallowtail/s.bitfinex/dto"
	"swallowtail/s.bitfinex/marshaling"
	bitfinexproto "swallowtail/s.bitfinex/proto"
)

// GetFundingRate fetches the funding rates from Bitfinex by `symbol` & `limit`.
func (s *BitfinexService) GetFundingRate(
	ctx context.Context, in *bitfinexproto.GetBitfinexFundingRatesRequest,
) (*bitfinexproto.GetBitfinexFundingRatesResponse, error) {
	// Validation.
	switch {
	case in.Symbol == "":
		return nil, gerrors.BadParam("missing_param.symbol", nil)
	}

	errParams := map[string]string{
		"symbol": in.Symbol,
		"limit":  strconv.Itoa(int(in.Limit)),
	}

	// Get funding rate.
	rsp, err := client.GetFundingRates(ctx, &dto.GetFundingRatesRequest{})
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_get_funding_rates", errParams)
	}

	// Marshal to proto.
	proto := marshaling.GetFundingRatesDTOToProto(rsp)

	return proto, nil
}
