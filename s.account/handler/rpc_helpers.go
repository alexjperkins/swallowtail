package handler

import (
	"context"
	"fmt"
	"time"

	"swallowtail/libraries/gerrors"
	accountproto "swallowtail/s.account/proto"
	binanceproto "swallowtail/s.binance/proto"
	discordproto "swallowtail/s.discord/proto"
)

func validateExchangeCredentials(ctx context.Context, userID string, exchange *accountproto.Exchange) (bool, string, error) {
	errParams := map[string]string{
		"exchange_type": exchange.ExchangeType.String(),
	}

	switch exchange.ExchangeType.String() {
	case accountproto.ExchangeType_BINANCE.String():
		return validateBinanceExchangeCredentials(ctx, userID, exchange)
	case accountproto.ExchangeType_FTX.String():
		return validateFTXExchangeCredentials(ctx, userID, exchange)
	default:
		return false, "", gerrors.FailedPrecondition("failed_to_validate_credentials.invalid_exchange", errParams)
	}
}

func validateBinanceExchangeCredentials(ctx context.Context, userID string, exchange *accountproto.Exchange) (bool, string, error) {
	rsp, err := (&binanceproto.VerifyCredentialsRequest{
		UserId: userID,
		Credentials: &binanceproto.Credentials{
			ApiKey:    exchange.ApiKey,
			SecretKey: exchange.SecretKey,
		},
	}).SendWithTimeout(ctx, 30*time.Second).Response()
	if err != nil {
		return false, "", gerrors.Augment(err, "failed_to_validate_binance_credentials", nil)
	}

	return rsp.Success, rsp.Reason, nil
}

func validateFTXExchangeCredentials(ctx context.Context, userID string, exchange *accountproto.Exchange) (bool, string, error) {
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
