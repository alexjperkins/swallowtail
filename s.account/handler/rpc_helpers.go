package handler

import (
	"context"
	"fmt"
	"time"

	"swallowtail/libraries/gerrors"
	accountproto "swallowtail/s.account/proto"
	binanceproto "swallowtail/s.binance/proto"
	discordproto "swallowtail/s.discord/proto"
	tradeengineproto "swallowtail/s.trade-engine/proto"
)

func validateVenueCredentials(ctx context.Context, userID string, venueAccount *accountproto.VenueAccount) (bool, string, error) {
	errParams := map[string]string{
		"venue": venueAccount.Venue.String(),
	}

	// Validate venue credentials.
	switch venueAccount.Venue {
	case tradeengineproto.VENUE_BINANCE:
		return validateBinanceCredentials(ctx, userID, venueAccount)
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

func validateBinanceCredentials(ctx context.Context, userID string, venueAccount *accountproto.VenueAccount) (bool, string, error) {
	rsp, err := (&binanceproto.VerifyCredentialsRequest{
		UserId: userID,
		Credentials: &tradeengineproto.VenueCredentials{
			ApiKey:    venueAccount.ApiKey,
			SecretKey: venueAccount.SecretKey,
		},
	}).SendWithTimeout(ctx, 30*time.Second).Response()
	if err != nil {
		return false, "", gerrors.Augment(err, "failed_to_validate_binance_credentials", nil)
	}

	return rsp.Success, rsp.Reason, nil
}

func validateFTXCredentials(ctx context.Context, userID string, venueAccount *accountproto.VenueAccount) (bool, string, error) {
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
