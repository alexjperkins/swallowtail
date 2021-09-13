package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/monzo/slog"

	"swallowtail/libraries/gerrors"
	discordproto "swallowtail/s.discord/proto"
)

func executeTradeForUser(ctx context.Context, userID, tradeID string, riskPercentage int) error {
	if err := notifyUserOnSuccess(ctx, userID, tradeID, riskPercentage); err != nil {
		slog.Error(ctx, "Failed to notify user of trade.", nil)
	}

	return nil
}

func notifyUserOnSuccess(ctx context.Context, userID, tradeID string, risk int) error {
	header := fmt.Sprintf(":wave: <@%s>, I have place your Trade with %s%% risk :rocket:", userID, risk)

	content := `
TRADE ID:  %s
ASSET:     %s
RISK (%%):%s
EXCHANGE:  %s
`
	formattedContent := fmt.Sprintf(content, tradeID, "", risk, "")

	if _, err := (&discordproto.SendMsgToPrivateChannelRequest{
		UserId:         userID,
		Content:        fmt.Sprintf("%s```%s```", header, formattedContent),
		IdempotencyKey: fmt.Sprintf("%s-%s-%s", userID, tradeID, time.Now().UTC().Truncate(15*time.Minute)),
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

func notifyPulseChannelUserTrade(ctx context.Context, userID, tradeID string, risk int) error {
	now := time.Now()

	header := ":dove:   `CRONITOR: TRADE PARTICIPANTS NEW PARTICIPANT`   :money_mouth:"
	content := `
TRADE ID:  %s
USER ID:   %s
TIMESTAMP: %v
RISK (%):  %v%%
`
	formattedContent := fmt.Sprintf(content, tradeID, userID, time.Now().UTC().Truncate(time.Second), risk)

	if _, err := (&discordproto.SendMsgToChannelRequest{
		ChannelId:      discordproto.DiscordSatoshiTradesPulseChannel,
		Content:        fmt.Sprintf("%s```%s```", header, formattedContent),
		IdempotencyKey: fmt.Sprintf("tradeparticipantspollnew-%s-%v", userID, now.Truncate(time.Hour)),
	}).Send(ctx).Response(); err != nil {
		return gerrors.Augment(err, "failed_to_notify_user", nil)
	}

	return nil
}
