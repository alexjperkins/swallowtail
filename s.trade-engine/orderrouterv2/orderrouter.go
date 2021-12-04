package or

import (
	"context"
	"strings"
	"swallowtail/libraries/gerrors"

	"github.com/hashicorp/go-multierror"
	"google.golang.org/protobuf/types/known/timestamppb"

	binanceproto "swallowtail/s.binance/proto"
	tradeengineproto "swallowtail/s.trade-engine/proto"
)

// RouteExecuteNewOrder ...
func RouteExecuteNewOrder(
	ctx context.Context,
	orders []*tradeengineproto.Order,
	venue tradeengineproto.VENUE,
	instrumentType tradeengineproto.INSTRUMENT_TYPE,
	credentials *tradeengineproto.VenueCredentials,
) ([]*tradeengineproto.Order, error) {
	errParams := map[string]string{
		"venue_id":        strings.ToLower(venue.String()),
		"instrument_type": strings.ToLower(instrumentType.String()),
	}

	switch venue {
	case tradeengineproto.VENUE_BINANCE:
		switch instrumentType {
		case tradeengineproto.INSTRUMENT_TYPE_FUTURE_PERPETUAL:
			return executeBinanceNewPerpetualFuturesOrders(ctx, orders, venueCredentials)
		case tradeengineproto.INSTRUMENT_TYPE_SPOT:
			return executeBinanceNewSpotOrders(ctx, orders, venueCredentials)
		default:
			return nil, gerrors.Unimplemented("failed_to_route_and_execute_order.instrument_exchange_pair_umimplemented", errParams)
		}
	case tradeengineproto.VENUE_FTX:
		return executeFTXNewOrders(ctx, orders, venueCredentials)
	default:
		return nil, gerrors.Unimplemented("failed_to_route_and_execute_order.exchange_unimplemented", errParams)
	}

	return nil, nil
}

// executeBinanceNewPerpetualFuturesOrder ...
func executeBinanceNewPerpetualFuturesOrders(ctx context.Context, orders []*tradeengineproto.Order, credentials *tradeengineproto.VenueCredentials) ([]*tradeengineproto.Order, error) {
	var (
		os   []*tradeengineproto.Order
		mErr error
	)
	for _, o := range orders {
		rsp, err := (&binanceproto.ExecuteNewFuturesPerpetualOrderRequest{
			Order:       o, // should compiler ignore this error?
			Credentials: credentials,
			Timestamp:   timestamppb.Now(),
		}).Send(ctx).Response()
		if err != nil {
			mErr = multierror.Append(mErr, err)
		}

		os = append(os, rsp.Order)
	}

	if mErr != nil {
		return os, gerrors.Augment(mErr, "failed_to_route_and_execute_order.binance_perpetual_futures", nil)
	}

	return os, nil
}

// executeBinanceNewSpotOrder ...
func executeBinanceNewSpotOrders(ctx context.Context, orders []*tradeengineproto.Order, credentials *tradeengineproto.VenueCredentials) ([]*tradeengineproto.Order, error) {
	return nil, gerrors.Unimplemented("failed_to_route_and_execute_order.ftx_new_orders_spot_unimplemented", nil)
}

// executeFTXNewOrder ...
func executeFTXNewOrders(ctx context.Context, orders []*tradeengineproto.Order, credentials *tradeengineproto.VenueCredentials) ([]*tradeengineproto.Order, error) {
	return nil, gerrors.Unimplemented("failed_to_route_and_execute_order.ftx_new_orders_unimplemented", nil)
}
