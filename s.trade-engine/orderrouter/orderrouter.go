package orderrouter

import (
	"context"
	"strings"

	"swallowtail/libraries/gerrors"
	tradeengineproto "swallowtail/s.trade-engine/proto"

	"github.com/monzo/slog"
)

// RouteAndExecuteNewOrder routes order to the correct exchange & executes.
func RouteAndExecuteNewOrder(
	ctx context.Context,
	order *tradeengineproto.Order,
	venue tradeengineproto.VENUE,
	instrumentType tradeengineproto.INSTRUMENT_TYPE,
	venueCredentials *tradeengineproto.VenueCredentials,
) (*tradeengineproto.Order, error) {

	errParams := map[string]string{
		"venue_id":        strings.ToLower(venue.String()),
		"instrument_type": strings.ToLower(instrumentType.String()),
	}

	switch venue {
	case tradeengineproto.VENUE_BINANCE:
		switch instrumentType {
		case tradeengineproto.INSTRUMENT_TYPE_FUTURE_PERPETUAL:
			return executeBinanceNewPerpetualFuturesOrders(ctx, order, venueCredentials)
		case tradeengineproto.INSTRUMENT_TYPE_SPOT:
			return executeBinanceNewSpotOrders(ctx, order, venueCredentials)
		default:
			return nil, gerrors.Unimplemented("failed_to_route_and_execute_order.instument_not_supported_on_venue", errParams)
		}
	case tradeengineproto.VENUE_FTX:
		switch instrumentType {
		case tradeengineproto.INSTRUMENT_TYPE_FORWARD:
			return nil, gerrors.Unimplemented("failed_to_route_and_execute_order.instument_not_supported_on_venue", errParams)
		default:
			return executeFTXNewOrders(ctx, order, venueCredentials)
		}
	default:
		slog.Error(ctx, "Failed to route order: venue, instrument pair not implemented: %+v", errParams)
		return nil, gerrors.Unimplemented("failed_to_route_and_execute_order.venue_unimplemented", errParams)
	}
}
