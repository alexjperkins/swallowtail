package handler

import (
	"context"
	"fmt"
	"time"

	"swallowtail/libraries/gerrors"
	discordproto "swallowtail/s.discord/proto"
)

func notifyUserOnFailure(ctx context.Context, userID, tradeID string, err error) error {
	var errMsg string
	switch {
	case gerrors.Is(err, gerrors.ErrUnauthenticated):
		errMsg = "You are unauthenticated; this means you likely have issues with your API keys."
	case gerrors.Is(err, gerrors.ErrFailedPrecondition, "failed_to_read_primary_exchange.account_required"):
		errMsg = "You must register an account before placing trades. See `!account register help`."
	case gerrors.Is(err, gerrors.ErrNotFound, "exchanges_not_found_for_user_id"):
		errMsg = "You don't have any exchange registered. `See !exchange register help`."
	case gerrors.Is(err, gerrors.ErrFailedPrecondition, "exchange_found_different_to_primary_exchange_on_account"):
		errMsg = "There is an issue with your primary exchange. Sorry about that; please ping @ajperkins for a hand."
	case gerrors.Is(err, gerrors.ErrAlreadyExists, "failed_to_add_participant_to_trade.trade_already_exists"):
		errMsg = "I already have a record of this trade. You've already placed it. If this is incorrect please ping @ajperkins."
	case gerrors.Is(err, gerrors.ErrUnknown):
		errMsg = "Sorry I'm not quite sure what happened. Please ping @ajperkins to investigate."
	case gerrors.Is(err, gerrors.ErrFailedPrecondition, "bad_request"):
		errMsg = "Sorry, the request to the exchange was malformed. This shouldnt' happen, **please** ping @ajperkins to investigate."
	case gerrors.Is(err, gerrors.ErrRateLimited):
		errMsg = "Sorry, looks like I've been rate limited. Please try and place the trade manually again in a few seconds time."
	default:
		errMsg = "Sorry, I'm not sure what happened there. Please ping @ajperkins for a hand."
	}

	header := fmt.Sprintf(":warning: <@%s>, I failed to fully place your Trade. Please manually check on the exchange :warning:\nIf the error is transient you can try to place manually with a command.", userID)
	content := `
TRADE ID: %s
ERROR:    %v
`
	formattedContent := fmt.Sprintf(content, tradeID, errMsg)

	if _, err := (&discordproto.SendMsgToPrivateChannelRequest{
		UserId:         userID,
		Content:        fmt.Sprintf("%s```%s```", header, formattedContent),
		IdempotencyKey: fmt.Sprintf("tradefailure-%s-%s-%s", userID, tradeID, time.Now().UTC().Truncate(15*time.Minute)),
	}).Send(ctx).Response(); err != nil {
		return gerrors.Augment(err, "failed_to_notify_user", nil)
	}

	return nil
}

func notifyUserOnSuccess(ctx context.Context, userID, tradeID, exchangeTradeID, tradeParticipantID, asset, exchange string, risk, size float64, timestamp time.Time) error {
	header := fmt.Sprintf(":wave: <@%s>, I have place your Trade with %v%% risk :rocket:", userID, risk)

	content := `
TRADE ID:             %s
EXCHANGE TRADE ID:    %s
TRADE PARTICIPANT ID: %s
ASSET:                %s
EXCHANGE:             %s
RISK (%%):            %v
SIZE:                 %v
TIMESTAMP:            %v
`
	formattedContent := fmt.Sprintf(content, tradeID, exchangeTradeID, tradeParticipantID, asset, exchange, risk, size, timestamp)

	if _, err := (&discordproto.SendMsgToPrivateChannelRequest{
		UserId:         userID,
		Content:        fmt.Sprintf("%s```%s```", header, formattedContent),
		IdempotencyKey: fmt.Sprintf("tradesuccess-%s-%s-%s", userID, tradeID, time.Now().UTC().Truncate(15*time.Minute)),
	}).Send(ctx).Response(); err != nil {
		return gerrors.Augment(err, "failed_to_notify_user", nil)
	}

	return nil
}

func notifyTradesChannelContextEnded(ctx context.Context, tradeID string) error {
	now := time.Now().UTC()

	header := ":octopus:   `TRADE CONTEXT ENDED`   :four_leaf_clover:"
	content := `
TRADE ID:  %s
TIMESTAMP: %s

The 15 minute context for this trade has now ended. If you still would like to place the trade, you can place manually with a command. Good Luck!

!trade <trade_id> <risk>
`
	formattedContent := fmt.Sprintf(content, tradeID, now)

	if _, err := (&discordproto.SendMsgToChannelRequest{
		ChannelId:      discordproto.DiscordSatoshiModTradesChannel,
		Content:        fmt.Sprintf("%s```%s```", header, formattedContent),
		IdempotencyKey: fmt.Sprintf("tradeparticipantspolltradectxend-%s-%v", tradeID, time.Now().UTC().Truncate(time.Minute)),
	}).Send(ctx).Response(); err != nil {
		return gerrors.Augment(err, "failed_to_notify_user", nil)
	}

	return nil
}

func notifyPulseChannelHeartbeat(ctx context.Context, tradeID string, deadline time.Time) error {
	now := time.Now()

	header := ":heartpulse:   `CRONITOR: TRADE PARTICIPANTS POLL PULSE`   :heartpulse:"
	content := `
TRADE ID:  %s
TIMESTAMP: %v
DEADLINE:  %v
`
	formattedContent := fmt.Sprintf(content, tradeID, now, deadline)

	if _, err := (&discordproto.SendMsgToChannelRequest{
		ChannelId:      discordproto.DiscordSatoshiTradesPulseChannel,
		Content:        fmt.Sprintf("%s```%s```", header, formattedContent),
		IdempotencyKey: fmt.Sprintf("tradeparticipantspollstart-%s-%v", tradeID, time.Now().UTC().Truncate(time.Minute)),
	}).Send(ctx).Response(); err != nil {
		return gerrors.Augment(err, "failed_to_notify_user", nil)
	}

	return nil
}

func notifyPulseChannelStart(ctx context.Context, tradeID string, deadline time.Time) error {
	now := time.Now()

	header := ":robot:   `CRONITOR: TRADE PARTICIPANTS CONTEXT STARTED`   :heartpulse:"
	content := `
TRADE ID:  %s
TIMESTAMP: %v
DEADLINE:  %v
`
	formattedContent := fmt.Sprintf(content, tradeID, now, deadline)

	if _, err := (&discordproto.SendMsgToChannelRequest{
		ChannelId:      discordproto.DiscordSatoshiTradesPulseChannel,
		Content:        fmt.Sprintf("%s```%s```", header, formattedContent),
		IdempotencyKey: fmt.Sprintf("tradeparticipantspollstart-%s", tradeID),
	}).Send(ctx).Response(); err != nil {
		return gerrors.Augment(err, "failed_to_notify_user", nil)
	}

	return nil
}

func notifyPulseChannelEnd(ctx context.Context, tradeID string, deadline time.Time) error {
	now := time.Now().UTC()

	header := ":robot:   `CRONITOR: TRADE PARTICIPANTS CONTEXT FINISHED`   :skull:"
	content := `
TRADE ID:  %s
TIMESTAMP: %v
DEADLINE:  %v
`
	formattedContent := fmt.Sprintf(content, tradeID, now, deadline)

	if _, err := (&discordproto.SendMsgToChannelRequest{
		ChannelId:      discordproto.DiscordSatoshiTradesPulseChannel,
		Content:        fmt.Sprintf("%s```%s```", header, formattedContent),
		IdempotencyKey: fmt.Sprintf("tradeparticipantspollend-%s", tradeID),
	}).Send(ctx).Response(); err != nil {
		return gerrors.Augment(err, "failed_to_notify_user", nil)
	}

	return nil
}

func notifyPulseChannelUserTradeSuccess(ctx context.Context, userID, tradeID string, risk int) error {
	now := time.Now()

	header := ":dove:   `CRONITOR: TRADE PARTICIPANTS NEW PARTICIPANT`   :money_mouth:"
	content := `
TRADE ID:  %s
USER ID:   %s
TIMESTAMP: %v
RISK (%%):  %v
`
	formattedContent := fmt.Sprintf(content, tradeID, userID, time.Now().UTC().Truncate(time.Second), risk)

	if _, err := (&discordproto.SendMsgToChannelRequest{
		ChannelId:      discordproto.DiscordSatoshiTradesPulseChannel,
		Content:        fmt.Sprintf("%s```%s```", header, formattedContent),
		IdempotencyKey: fmt.Sprintf("tradeparticipantspollnewsuccess-%s-%v", userID, now.Truncate(time.Hour)),
	}).Send(ctx).Response(); err != nil {
		return gerrors.Augment(err, "failed_to_notify_user", nil)
	}

	return nil
}

func notifyPulseChannelUserTradeFailure(ctx context.Context, userID, tradeID string, risk int, err error) error {
	now := time.Now()

	header := ":rotating_light:   `CRONITOR: TRADE PARTICIPANTS TRADE FAILED`   :warning:"
	content := `
TRADE ID:  %s
USER ID:   %s
TIMESTAMP: %v
RISK (%%): %v

ERROR:     %v
`
	formattedContent := fmt.Sprintf(content, tradeID, userID, time.Now().UTC().Truncate(time.Second), risk, err)

	if _, err := (&discordproto.SendMsgToChannelRequest{
		ChannelId:      discordproto.DiscordSatoshiTradesPulseChannel,
		Content:        fmt.Sprintf("%s```%s```", header, formattedContent),
		IdempotencyKey: fmt.Sprintf("tradeparticipantspollnewfailure-%s-%v", userID, now.Truncate(time.Hour)),
	}).Send(ctx).Response(); err != nil {
		return gerrors.Augment(err, "failed_to_notify_user", nil)
	}

	return nil
}
