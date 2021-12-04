package execution

import (
	"context"
	"fmt"
	"time"

	"swallowtail/libraries/gerrors"
	"swallowtail/libraries/util"
	accountproto "swallowtail/s.account/proto"
	binanceproto "swallowtail/s.binance/proto"
	discordproto "swallowtail/s.discord/proto"
	"swallowtail/s.trade-engine/marshaling"
	tradeengineproto "swallowtail/s.trade-engine/proto"
)

func readVenueCredentials(ctx context.Context, userID string, venue tradeengineproto.VENUE) (*tradeengineproto.VenueCredentials, error) {
	rsp, err := (&accountproto.ReadExchangeByExchangeDetailsRequest{
		Exchange: venue.String(),
		UserId:   userID,
	}).Send(ctx).Response()
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_read_exchange_credentials", nil)
	}

	// Marshal venue credentials.
	venueCredentials, err := marshaling.AccountExchangeToVenueCredentials(rsp.GetExchange())
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_route_and_execute_order.marshal_credentials", nil)
	}

	return venueCredentials, nil
}

func readVenueAccountBalance(ctx context.Context, venue tradeengineproto.VENUE, _ tradeengineproto.INSTRUMENT_TYPE, credentials *tradeengineproto.VenueCredentials) (float64, error) {
	errParams := map[string]string{
		"venue": venue.String(),
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

func notifyUser(ctx context.Context, msg string, userID string) error {
	content := `:wave: <@%s> WARNING FROM TRADE ENGINE:

%s
`

	formattedContent := fmt.Sprintf(content, userID, msg)
	if _, err := (&discordproto.SendMsgToPrivateChannelRequest{
		UserId:         userID,
		SenderId:       tradeengineproto.TradeEngineActorSatoshiSystem,
		Content:        formattedContent,
		IdempotencyKey: fmt.Sprintf("%s-%s-%s", userID, util.Sha256Hash(msg), time.Now().UTC().Truncate(10*time.Minute)),
	}).Send(ctx).Response(); err != nil {
		return gerrors.Augment(err, "failed_to_notify_user", nil)
	}

	return nil
}
