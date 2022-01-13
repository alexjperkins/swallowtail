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
	// Read venue account.
	rsp, err := (&accountproto.ReadVenueAccountByVenueAccountDetailsRequest{
		Venue:          venue,
		UserId:         userID,
		ActorId:        accountproto.ActorSystemTradeEngine,
		RequestContext: accountproto.RequestContextOrderRequest,
	}).Send(ctx).Response()
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_read_venue_credentials", nil)
	}

	// Marshal venue credentials.
	return marshaling.VenueAccountToVenueCredentials(rsp.GetVenueAccount()), nil
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
