package handler

import (
	"context"
	"fmt"
	"time"

	"swallowtail/libraries/gerrors"
	discordproto "swallowtail/s.discord/proto"
)

func notifyUserOnFailure(ctx context.Context, userID, tradeID string, err error) error {
	header := fmt.Sprintf(":warning: <@%s>, I failed to fully place your Trade. Please manually check on the exchange :warning:\nIf the error is transient you can try to place manually with a command.", userID)
	content := `
TRADE ID: %s
ERROR:    %v
`
	formattedContent := fmt.Sprintf(content, tradeID, err)

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
RISK (%):             %v
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
RISK (%):  %v
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
RISK (%):  %v

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
