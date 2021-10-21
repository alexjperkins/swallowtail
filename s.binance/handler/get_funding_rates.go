package handler

import (
	"context"
	"strconv"
	"swallowtail/libraries/gerrors"
	"swallowtail/s.binance/client"
	binanceproto "swallowtail/s.binance/proto"
)

// GetFundingRates ...
func (*BinanceService) GetFundingRates(
	ctx context.Context, in *binanceproto.GetFundingRateRequest,
) (*binanceproto.GetFundingRateResponse, error) {
	switch {
	case in.Symbol == "":
		return nil, gerrors.BadParam("missing_param.symbols", nil)
	}

	errParams := map[string]string{
		"symbol":     in.Symbol,
		"start_time": strconv.Itoa(int(in.StartTime)),
		"end_time":   strconv.Itoa(int(in.StartTime)),
	}

	rsp, err := client.GetFundingRate(ctx, &client.GetFundingRateRequest{
		Symbol:    in.Symbol,
		StartTime: int(in.StartTime),
		EndTime:   int(in.EndTime),
	})
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_get_funding_rates", errParams)
	}

	protos := make([]*binanceproto.FundingRateInfo, 0, len(*rsp))
	for _, fr := range *rsp {
		protos = append(protos, &binanceproto.FundingRateInfo{
			Symbol:          fr.Symbol,
			FundingRate:     float32(fr.FundingRate),
			FundingRateTime: float32(fr.FundingTime),
		})
	}

	return &binanceproto.GetFundingRateResponse{
		FundingRates: protos,
	}, nil
}
