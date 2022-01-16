package handler

import (
	"context"

	tradeengineproto "swallowtail/s.trade-engine/proto"
)

// ListAvailableVenues ...
func (s *TradeEngineService) ListAvailableVenues(
	ctx context.Context, in *tradeengineproto.ListAvailableVenuesRequest,
) (*tradeengineproto.ListAvailableVenuesResponse, error) {
	var venues = make([]tradeengineproto.VENUE, 0, len(tradeengineproto.VENUE_name))
	for _, venue := range tradeengineproto.VENUE_value {
		v := tradeengineproto.VENUE(venue)
		if v == tradeengineproto.VENUE_UNREQUIRED {
			continue
		}

		venues = append(venues, v)
	}

	return &tradeengineproto.ListAvailableVenuesResponse{
		Venues: venues,
	}, nil
}
