package execution

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"google.golang.org/protobuf/internal/version"
	"google.golang.org/protobuf/types/known/timestamppb"

	"swallowtail/libraries/gerrors"
	accountproto "swallowtail/s.account/proto"
	binanceproto "swallowtail/s.binance/proto"
	ftxproto "swallowtail/s.ftx/proto"
	tradeengineproto "swallowtail/s.trade-engine/proto"
)

// RouteExecuteNewOrder ...
func RouteExecuteNewOrder(
	ctx context.Context,
	orders []*tradeengineproto.Order,
	venue tradeengineproto.VENUE,
	instrumentType tradeengineproto.INSTRUMENT_TYPE,
	credentials *accountproto.Exchange,
) ([]*tradeengineproto.Order, error) {
	errParams := map[string]string{
		"venue_id":        strings.ToLower(venue.String()),
		"instrument_type": strings.ToLower(instrumentType.String()),
	}

	switch venue {
	case tradeengineproto.VENUE_BINANCE:
		creds := &binanceproto.Credentials{
			ApiKey:    credentials.ApiKey,
			SecretKey: credentials.SecretKey,
		}

		switch instrumentType {
		case tradeengineproto.INSTRUMENT_TYPE_FUTURE_PERPETUAL:
			return executeBinanceNewPerpetualFuturesOrders(ctx, orders, creds)
		case tradeengineproto.INSTRUMENT_TYPE_SPOT:
			return executeBinanceNewSpotOrders(ctx, orders, creds)
		default:
			return nil, gerrors.Unimplemented("failed_to_route_and_execute_order.instrument_exchange_pair_umimplemented", errParams)
		}
	case tradeengineproto.VENUE_FTX:
		creds := &ftxproto.FTXCredentials{
			ApiKey:     credentials.ApiKey,
			SecretKey:  credentials.SecretKey,
			Subaccount: credentials.SubAccount,
		}

		return executeFTXNewOrders(ctx, orders, creds)
	default:
		return nil, gerrors.Unimplemented("failed_to_route_and_execute_order.exchange_unimplemented", errParams)
	}

	return nil, nil
}

// executeBinanceNewPerpetualFuturesOrder ...
func executeBinanceNewPerpetualFuturesOrders(ctx context.Context, orders []*tradeengineproto.Order, credentials *binanceproto.Credentials) ([]*tradeengineproto.Order, error) {
	var (
		os   []*tradeengineproto.Order
		mErr error
	)
	for _, o := range orders {
		rsp, err := (&binanceproto.ExecuteNewFuturesPerpetualOrderRequest{
			Order:       o,
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

func readVenueCredentials(ctx context.Context, userID string, venue tradeengineproto.VENUE) (*tradeengineproto.VenueCredentials, error) {
	rsp, err := (&accountproto.ReadExchangeByExchangeDetailsRequest{
		Exchange: venue.String(),
		UserId:   userID,
	}).Send(ctx).Response()
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_read_exchange_credentials", nil)
	}

	e := rsp.GetExchange()
	return &tradeengineproto.VenueCredentials{
		ApiKey:     e.ApiKey,
		SecretKey:  e.SecretKey,
		Subaccount: e.SubAccount,
		Passphrase: "", // TODO: we first need to add this to the data model of an exchange.
	}, nil
}

func readVenueAccountBalance(ctx context.Context, venue tradeengineproto.VENUE, credentials *tradeengineproto.VenueCredentials) (float64, error) {
	errParams := map[string]string{
		"venue": version.String(),
	}

	switch venue {
	case tradeengineproto.VENUE_BINANCE:
		rsp, err := (&binanceproto.ReadPerpetualFuturesAccountRequest{
			ActorId:     binanceproto.BinanceAccountActorTradeEngineSystem,
			Credentials: credentials,
		}).Send(ctx).Response()
		if err != nil {
			return 0, gerrors.Augment(err, "failed_to_read_venue_account_balance", errParams)
		}

		return float64(rsp.Balance), nil
	case tradeengineproto.VENUE_BITFINEX:
		return 0, gerrors.Unimplemented("failed_to_read_venue_account_balance.unimplemented.venue", errParams)
	case tradeengineproto.VENUE_DERIBIT:
		return 0, gerrors.Unimplemented("failed_to_read_venue_account_balance.unimplemented.venue", errParams)
	case tradeengineproto.VENUE_FTX:
		return 0, gerrors.Unimplemented("failed_to_read_venue_account_balance.unimplemented.venue", errParams)
	default:
		return 0, gerrors.Unimplemented("failed_to_read_venue_account_balance.unimplemented.venue", errParams)
	}
}
