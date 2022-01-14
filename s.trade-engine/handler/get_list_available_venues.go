package handler

import (
	"context"
	"strings"

	tradeengineproto "swallowtail/s.trade-engine/proto"
)

// ListAvailableVenues ...
func (s *TradeEngineService) ListAvailableVenues(
	ctx context.Context, in *tradeengineproto.ListAvailableVenuesRequest,
) (*tradeengineproto.ListAvailableVenuesResponse, error) {
	var venues = make([]string, 0, len(tradeengineproto.VENUE_name))
	for _, v := range tradeengineproto.VENUE_name {
		if strings.ToUpper(v) != tradeengineproto.VENUE_UNREQUIRED.String() {
			continue
		}

		venues = append(venues, strings.ToUpper(v))
	}

	return &tradeengineproto.ListAvailableVenuesResponse{
		Venues: venues,
	}, nil
}
