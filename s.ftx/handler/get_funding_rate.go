package handler

import (
	"context"
	"swallowtail/libraries/gerrors"
	"swallowtail/s.ftx/client"
	ftxproto "swallowtail/s.ftx/proto"
)

// GetFundingRate ...
func (s *FTXService) GetFundingRate(
	ctx context.Context,
	in *ftxproto.GetFTXFundingRatesRequest,
) (*ftxproto.GetFTXFundingRatesResponse, error) {
	switch {
	case in.Symbol == "":
		return nil, gerrors.BadParam("missing_param.symbol", nil)
	}

	errParams := map[string]string{
		"symbol": in.Symbol,
	}

	var limit int = 1
	if in.Limit > 1 {
		limit = int(in.Limit)
	}

	rsp, err := client.GetFundingRate(ctx, &client.GetFundingRateRequest{
		Instrument: in.Symbol,
		StartTime:  int(in.StartTime),
		EndTime:    int(in.EndTime),
		Limit:      limit,
	})
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_get_funding_rated", errParams)
	}

	protos := make([]*ftxproto.FTXFundingRatesInfo, 0, len(rsp.FundingRates))
	for _, fr := range rsp.FundingRates {
		protos = append(protos, &ftxproto.FTXFundingRatesInfo{
			Symbol:      fr.Instrument,
			FundingRate: float32(fr.Rate),
			FundingTime: fr.Time,
		})
	}

	return &ftxproto.GetFTXFundingRatesResponse{
		FundingRates: protos,
	}, nil
}
