package handler

import (
	"context"
	"fmt"
	"time"

	"swallowtail/libraries/gerrors"
	discordproto "swallowtail/s.discord/proto"
	tradeengineproto "swallowtail/s.trade-engine/proto"
)

func notifyUserOnFailure(ctx context.Context, userID, tradeStrategyID string, numberOfSuccessOrders int, err error, executionError *tradeengineproto.ExecutionError) error {
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
		errMsg = "Sorry, the request to the exchange was malformed. This can happen if the trade amount you place is too small, **please** ping @ajperkins to investigate if you don't believe this is the case."
	case gerrors.Is(err, gerrors.ErrRateLimited):
		errMsg = "Sorry, looks like I've been rate limited. Please try and place the trade manually again in a few seconds time."
	default:
		errMsg = "Sorry, I'm not sure what happened there. Please ping @ajperkins for a hand."
	}

	header := fmt.Sprintf(
		":warning: <@%s>, I failed to fully execute your trade strategy, %d were placed. Please manually check on the exchange :warning:\nIf the error is transient you may try to place manually with a command.",
		userID,
		numberOfSuccessOrders,
	)

	content := `
TRADE STRATEGY ID: %s
ERROR:             %v
EXECUTION ERROR:   %v
FAILED ORDER:      %+v
`
	formattedContent := fmt.Sprintf(content, tradeStrategyID, errMsg, executionError.GetErrorMessage(), executionError.GetFailedOrder())

	if _, err := (&discordproto.SendMsgToPrivateChannelRequest{
		UserId:         userID,
		Content:        fmt.Sprintf("%s```%s```", header, formattedContent),
		IdempotencyKey: fmt.Sprintf("tradestrategyfailure-%s-%s-%s", userID, tradeStrategyID, time.Now().UTC().Truncate(15*time.Minute)),
	}).Send(ctx).Response(); err != nil {
		return gerrors.Augment(err, "failed_to_notify_user", nil)
	}

	return nil
}

func notifyUserOnSuccess(
	ctx context.Context,
	userID, tradeStrategyID, tradeParticipantID, asset, pair string,
	executionStrategy tradeengineproto.EXECUTION_STRATEGY,
	venue tradeengineproto.VENUE,
	risk, size float64,
	timestamp time.Time,
	successfulOrders []*tradeengineproto.Order,
) error {
	header := fmt.Sprintf(":wave: <@%s>, I have executed your trade strategy with %v%% risk :rocket:", userID, risk)

	content := `
TRADE STRATEGY ID:    %s
TRADE PARTICIPANT ID: %s
ASSET:                %s
PAIR:                 %s
VENUE:                %s
EXECUTION STRATEGY:   %v
RISK (%%):             %v
SIZE:                 %v
TIMESTAMP:            %v

SUCCESSFUL_ORDERS:    %v
`
	formattedContent := fmt.Sprintf(content, tradeStrategyID, tradeParticipantID, asset, pair, venue, executionStrategy, risk, size, timestamp, len(successfulOrders))

	if _, err := (&discordproto.SendMsgToPrivateChannelRequest{
		UserId:         userID,
		Content:        fmt.Sprintf("%s```%s```", header, formattedContent),
		IdempotencyKey: fmt.Sprintf("tradestrategysuccess-%s-%s-%s", userID, tradeStrategyID, time.Now().UTC().Truncate(15*time.Minute)),
	}).Send(ctx).Response(); err != nil {
		return gerrors.Augment(err, "failed_to_notify_user", nil)
	}

	return nil
}

func notifyTradesChannelContextEnded(ctx context.Context, tradeID string) error {
	now := time.Now().UTC()

	header := ":octopus:   `TRADE CONTEXT ENDED`   :four_leaf_clover:"
	content := `
TRADE STRATEGY ID:  %s
TIMESTAMP:          %s

The 15 minute context for this trade strategy has now ended. If you still would like to execute the trade strategy, you can place manually with a command. Good Luck!

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

	header := ":heartpulse:   `CRONITOR: TRADE STRATEGY PARTICIPANTS POLL PULSE`   :heartpulse:"
	content := `
TRADE STRATEGY ID:  %s
TIMESTAMP:          %v
DEADLINE:           %v
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

	header := ":robot:   `CRONITOR: TRADE STRATEGY PARTICIPANTS CONTEXT STARTED`   :heartpulse:"
	content := `
TRADE STRATEGY ID:  %s
TIMESTAMP:          %v
DEADLINE:           %v
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

	header := ":robot:   `CRONITOR: TRADE STRATEGY PARTICIPANTS CONTEXT FINISHED`   :skull:"
	content := `
TRADE STRATEGY ID:  %s
TIMESTAMP:          %v
DEADLINE:           %v
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

func notifyPulseChannelUserTradeSuccess(ctx context.Context, userID, tradeID string, executionStrategy tradeengineproto.EXECUTION_STRATEGY, venue tradeengineproto.VENUE, risk int, succesfulOrders []*tradeengineproto.Order) error {
	now := time.Now().UTC()

	header := ":dove:   `CRONITOR: TRADE STRATEGY NEW PARTICIPANT`   :money_mouth:"
	content := `
TRADE STRATEGY ID:  %s
USER ID:            %s
TIMESTAMP:          %v
EXECUTION STRATEGY: %v
VENUE:              %v
RISK (%%):           %v
SUCCESSFUL ORDERS:  %v
`
	formattedContent := fmt.Sprintf(content, tradeID, userID, time.Now().UTC().Truncate(time.Second), executionStrategy, venue, risk, len(succesfulOrders))

	if _, err := (&discordproto.SendMsgToChannelRequest{
		ChannelId:      discordproto.DiscordSatoshiTradesPulseChannel,
		Content:        fmt.Sprintf("%s```%s```", header, formattedContent),
		IdempotencyKey: fmt.Sprintf("tradeparticipantspollnewsuccess-%s-%v", userID, now.Truncate(time.Hour)),
	}).Send(ctx).Response(); err != nil {
		return gerrors.Augment(err, "failed_to_notify_user", nil)
	}

	return nil
}

func notifyPulseChannelUserTradeFailure(ctx context.Context, userID, tradeID string, risk, numberOfSuccessOrders int, err error, executionError *tradeengineproto.ExecutionError) error {
	now := time.Now()

	header := ":rotating_light:   `CRONITOR: TRADE STRATEGY PARTICIPANT EXECUTION FAILED`   :warning:"
	content := `
TRADE STRATEGY ID:  %s
USER ID:            %s
TIMESTAMP:          %v
RISK (%%):          %v
SUCCESSFUL_ORDERS:  %d
ERROR:              %v
EXECUTION_ERROR:    %v
FAILED_ORDER:       %+v
`
	formattedContent := fmt.Sprintf(content, tradeID, userID, time.Now().UTC().Truncate(time.Second), risk, numberOfSuccessOrders, err, executionError.GetErrorMessage(), executionError.GetFailedOrder())

	if _, err := (&discordproto.SendMsgToChannelRequest{
		ChannelId:      discordproto.DiscordSatoshiTradesPulseChannel,
		Content:        fmt.Sprintf("%s```%s```", header, formattedContent),
		IdempotencyKey: fmt.Sprintf("tradeparticipantspollnewfailure-%s-%v", userID, now.Truncate(time.Hour)),
	}).Send(ctx).Response(); err != nil {
		return gerrors.Augment(err, "failed_to_notify_user", nil)
	}

	return nil
}
