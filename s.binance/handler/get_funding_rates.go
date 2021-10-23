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
	ctx context.Context, in *binanceproto.GetFundingRatesRequest,
) (*binanceproto.GetFundingRatesResponse, error) {
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
		Limit:     int(in.Limit),
	})
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_get_funding_rates", errParams)
	}

	protos := make([]*binanceproto.FundingRateInfo, 0, len(*rsp))
	for _, fri := range *rsp {
		fr, err := strconv.ParseFloat(fri.FundingRate, 32)
		if err != nil {
			return nil, gerrors.Augment(err, "failed_to_get_funding_rates", errParams)
		}
		protos = append(protos, &binanceproto.FundingRateInfo{
			Symbol:          fri.Symbol,
			FundingRate:     float32(fr),
			FundingRateTime: float32(fri.FundingTime),
		})
	}

	return &binanceproto.GetFundingRatesResponse{
		FundingRates: protos,
	}, nil
}
