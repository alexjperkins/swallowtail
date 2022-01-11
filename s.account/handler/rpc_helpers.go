package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/monzo/slog"

	"swallowtail/libraries/gerrors"
	"swallowtail/s.account/domain"
	"swallowtail/s.account/marshaling"
	binanceproto "swallowtail/s.binance/proto"
	discordproto "swallowtail/s.discord/proto"
	tradeengineproto "swallowtail/s.trade-engine/proto"
)

func validateVenueCredentials(ctx context.Context, userID string, venueAccount interface{}) (bool, string, error) {
	var credentials *tradeengineproto.VenueCredentials
	switch t := venueAccount.(type) {
	case *domain.VenueAccount:
		venue, err := marshaling.ConvertVenueIDToProto(t.VenueID)
		if err != nil {
			return false, "", gerrors.Augment(err, "failed_to_validate_credentials.venue_account", nil)
		}

		credentials = &tradeengineproto.VenueCredentials{
			Venue:      venue,
			ApiKey:     t.APIKey,
			SecretKey:  t.SecretKey,
			Subaccount: t.SubAccount,
			Url:        t.URL,
			WsUrl:      t.WSURL,
		}
	case *domain.InternalVenueAccount:
		venue, err := marshaling.ConvertVenueIDToProto(t.VenueID)
		if err != nil {
			return false, "", gerrors.Augment(err, "failed_to_validate_credentials.internal_venue_account", nil)
		}

		credentials = &tradeengineproto.VenueCredentials{
			Venue:      venue,
			ApiKey:     t.APIKey,
			SecretKey:  t.SecretKey,
			Subaccount: t.SubAccount,
			Url:        t.URL,
			WsUrl:      t.WSURL,
		}
	default:
		slog.Error(ctx, "Failed to validate venue credentials, invalid type: %T", t)
		return false, "", gerrors.Unimplemented("unimplemented_account_type", map[string]string{
			"account_type": fmt.Sprintf("%T", t),
		})
	}

	errParams := map[string]string{
		"venue": credentials.Venue.String(),
	}

	// Validate venue credentials.
	switch credentials.Venue {
	case tradeengineproto.VENUE_BINANCE:
		return validateBinanceCredentials(ctx, userID, credentials)
	case tradeengineproto.VENUE_BITFINEX:
		return false, "", gerrors.Unimplemented("venue_unimplemented.bitfinex", nil)
	case tradeengineproto.VENUE_DERIBIT:
		return false, "", gerrors.Unimplemented("venue_unimplemented.deribit", nil)
	case tradeengineproto.VENUE_FTX:
		return false, "", gerrors.Unimplemented("venue_validation_unimplemented.ftx", nil)
	default:
		return false, "", gerrors.FailedPrecondition("failed_to_validate_credentials.invalid_venue_account", errParams)
	}
}

func validateBinanceCredentials(ctx context.Context, userID string, venueCredentials *tradeengineproto.VenueCredentials) (bool, string, error) {
	rsp, err := (&binanceproto.VerifyCredentialsRequest{
		UserId:      userID,
		Credentials: venueCredentials,
	}).SendWithTimeout(ctx, 30*time.Second).Response()
	if err != nil {
		return false, "", gerrors.Augment(err, "failed_to_validate_binance_credentials", nil)
	}

	return rsp.Success, rsp.Reason, nil
}

func validateFTXCredentials(ctx context.Context, userID string, venueCredentials *tradeengineproto.VenueCredentials) (bool, string, error) {
	return false, "", nil
}

func notifyPulseChannel(ctx context.Context, userID, username string, timestamp time.Time) error {
	base := ":bear:    `NEW MEMBER`    :bear:"
	msg := `
UserID: %s
Username: %s
Timestamp: %v
`
	formattedMsg := fmt.Sprintf(msg, userID, username, timestamp)

	if _, err := (&discordproto.SendMsgToChannelRequest{
		ChannelId: discordproto.DiscordSatoshiAccountsPulseChannel,
		Content:   fmt.Sprintf("%s```%s```", base, formattedMsg),
	}).Send(ctx).Response(); err != nil {
		return gerrors.Augment(err, "failed_to_notify_account_pulse_channel", nil)
	}

	return nil
}
