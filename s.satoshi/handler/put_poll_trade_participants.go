package handler

import (
	"context"
	"strconv"
	"time"

	"github.com/monzo/slog"

	"swallowtail/libraries/emojis"
	"swallowtail/libraries/gerrors"
	discordproto "swallowtail/s.discord/proto"
	satoshiproto "swallowtail/s.satoshi/proto"
)

// PollTradeParticipants ...
func (s *SatoshiService) PollTradeParticipants(
	ctx context.Context, in *satoshiproto.PollTradeParticipantsRequest,
) (*satoshiproto.PollTradeParticipantsResponse, error) {
	switch {
	case in.ActorId == "":
		return nil, gerrors.BadParam("missing_param.actor_id", nil)
	case !validateActor(in.ActorId):
		return nil, gerrors.Unauthenticated("invalid_actor", nil)
	case in.MessageId == "":
		return nil, gerrors.BadParam("missing_param.message_id", nil)
	case in.TradeId == "":
		return nil, gerrors.BadParam("missing_param.trade_id", nil)
	case in.ChannelId == "":
		return nil, gerrors.BadParam("missing_param.channel_id", nil)
	case in.TimeoutInMinutes <= 0:
		return nil, gerrors.BadParam("invalid_param.timeout_in_seconds.must_be_greater_than_zero", nil)
	}

	errParams := map[string]string{
		"message_id":         in.MessageId,
		"actor_id":           in.ActorId,
		"timeout_in_seconds": strconv.Itoa(int(in.TimeoutInMinutes)),
		"trade_id":           in.TradeId,
		"channel_id":         in.TradeId,
	}

	// This is horrible code; but we don't yet have the infra in place to do any better.
	// Ideally this should be asyncronous using some message queue using exactly-once semantics.
	go func() {
		deadline := time.Now().UTC().Add(time.Duration(in.TimeoutInMinutes) * time.Minute)
		newCtx := context.Background()
		childCtx, cancel := context.WithDeadline(newCtx, deadline)
		defer cancel()

		tradeCache := map[string]bool{}
		t := time.NewTicker(10 * time.Second)
		tPulse := time.NewTicker(5 * time.Minute)

		// Cronitor; notify pulse channel poll has started.
		if err := notifyPulseChannelStart(childCtx, in.TradeId, deadline); err != nil {
			slog.Error(childCtx, err.Error())
		}

		for {
			select {
			case <-tPulse.C:
				// Cronitor; notify pulse channel poll of pulse.
				if err := notifyPulseChannelHeartbeat(childCtx, in.TradeId, deadline); err != nil {
					slog.Error(newCtx, err.Error())
				}
			case <-t.C:
				// Poll for new reactions.
				rsp, err := (&discordproto.ReadMessageReactionsRequest{
					MessageId: in.MessageId,
					ChannelId: in.ChannelId,
				}).Send(childCtx).Response()
				if err != nil {
					slog.Trace(childCtx, "poll_trade_participants.failed_to_read_message_reactions", errParams)
					continue
				}

				reactions := rsp.GetReactions()
				for _, reaction := range reactions {
					if !isValidTradeParticipantReaction(reaction.GetReactionId()) {
						continue
					}

					for _, userID := range reaction.UserIds {
						if _, exists := tradeCache[userID]; exists {
							continue
						}

						risk := emojis.SatoshiRiskEmoji(reaction.GetReactionId()).AsRiskPercentage()
						if err := executeTradeForUser(newCtx, userID, in.TradeId, risk); err != nil {
							slog.Error(childCtx, "Failed to execute trade for user: %s", userID)
							continue
						}

						tradeCache[userID] = true

						if err := notifyPulseChannelUserTrade(childCtx, userID, in.TradeId, risk); err != nil {
							slog.Error(newCtx, err.Error())
						}
					}
				}

			case <-childCtx.Done():
				slog.Warn(newCtx, "Closing window for new trade participants for trade: %v", in.TradeId)
				if err := notifyPulseChannelEnd(newCtx, in.TradeId, deadline); err != nil {
					slog.Error(newCtx, err.Error())
				}

				return
			}
		}
	}()

	return &satoshiproto.PollTradeParticipantsResponse{}, nil
}