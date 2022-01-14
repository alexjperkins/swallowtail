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

// PollTradeStrategyParticipants polls a trade strategy formatted as a discord message for participant reactions,
// of which then executes that trade strategy for said participant.
func (s *SatoshiService) PollTradeStrategyParticipants(
	ctx context.Context, in *satoshiproto.PollTradeStrategyParticipantsRequest,
) (*satoshiproto.PollTradeStrategyParticipantsResponse, error) {
	switch {
	case in.ActorId == "":
		return nil, gerrors.BadParam("missing_param.actor_id", nil)
	case !validateActor(in.ActorId):
		return nil, gerrors.Unauthenticated("invalid_actor", nil)
	case in.MessageId == "":
		return nil, gerrors.BadParam("missing_param.message_id", nil)
	case in.TradeStrategyId == "":
		return nil, gerrors.BadParam("missing_param.trade_strategy_id", nil)
	case in.ChannelId == "":
		return nil, gerrors.BadParam("missing_param.channel_id", nil)
	case in.TimeoutInMinutes <= 0:
		return nil, gerrors.BadParam("invalid_param.timeout_in_seconds.must_be_greater_than_zero", nil)
	}

	errParams := map[string]string{
		"message_id":         in.MessageId,
		"actor_id":           in.ActorId,
		"timeout_in_seconds": strconv.Itoa(int(in.TimeoutInMinutes)),
		"trade_strategy_id":  in.TradeStrategyId,
		"channel_id":         in.TradeStrategyId,
	}

	// This is horrible code; but we don't yet have the infra in place to do any better.
	// Ideally this should be asynchronous using some message queue using exactly-once semantics.
	go func() {
		deadline := time.Now().UTC().Add(time.Duration(in.TimeoutInMinutes) * time.Minute)

		// We create a new context object; otherwise the parent context would be cancel once the
		// the response is returned to the callee.
		newCtx := context.Background()
		childCtx, cancel := context.WithDeadline(newCtx, deadline)
		defer cancel()

		tradeCache := map[string]bool{}
		t := time.NewTicker(10 * time.Second)
		tPulse := time.NewTicker(5 * time.Minute)

		// Cronitor; notify pulse channel poll has started.
		if err := notifyPulseChannelStart(childCtx, in.TradeStrategyId, deadline); err != nil {
			slog.Error(childCtx, err.Error())
		}

		// Poll for trade participants; ideally this is a candidate to be refactored.
		for {
			select {
			case <-tPulse.C:
				// Cronitor; notify pulse channel poll of pulse.
				if err := notifyPulseChannelHeartbeat(childCtx, in.TradeStrategyId, deadline); err != nil {
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

				// Execute order for particicpant based on trade strategy reaction.
				reactions := rsp.GetReactions()
				for _, reaction := range reactions {
					if !isValidTradeParticipantReaction(reaction.GetReactionId()) {
						continue
					}

					for _, userID := range reaction.UserIds {
						if _, exists := tradeCache[userID]; exists {
							continue
						}

						// We set this to true, even if we fail; we may want to introduce some retry mechanics.
						// But lets keep it simple for now.
						tradeCache[userID] = true

						// Calculate risk based on emoji.
						risk := emojis.SatoshiRiskEmoji(reaction.GetReactionId()).AsRiskPercentage()

						// Execute order.
						rsp, err := executeTradeStrategyForParticipant(newCtx, userID, in.TradeStrategyId, risk)
						if err != nil {
							slog.Error(newCtx, "Failed to execute trade strategy for user: %s; Error: %v", userID, err)

							// Notify parties of failure.
							if perr := notifyUserOnFailure(newCtx, userID, in.TradeStrategyId, 0, err, nil); perr != nil {
								slog.Error(newCtx, "Failed to notify user of successful trade strategy: %s, UserID %s, Error: %s", in.TradeStrategyId, userID, perr)
							}

							if perr := notifyPulseChannelUserTradeFailure(newCtx, userID, in.TradeStrategyId, risk, 0, err, nil); perr != nil {
								slog.Error(newCtx, "Failed to notify channel of successful trade strategy: TradeID %s, UserID %s, Error: %v", in.TradeStrategyId, userID, perr)
							}

							continue
						}

						// Verify we don't have a partial failure.
						if rsp.GetError() != nil {
							slog.Error(newCtx, "Partially failed to execute trade strategy for participant: %s; Error: %v", userID, rsp.GetError())

							// Notify pulse channel of partial failure.
							if perr := notifyPulseChannelUserTradeFailure(newCtx, userID, in.TradeStrategyId, risk, int(rsp.NumberOfExecutedOrders), nil, rsp.GetError()); perr != nil {
								slog.Error(newCtx, "Failed to notify channel of successful trade strategy: %s, UserID %s, Error: %v", in.TradeStrategyId, userID, perr)
							}
						}

						// Notify parties of success.
						if err := notifyUserOnSuccess(
							newCtx,
							userID,
							in.TradeStrategyId,
							userID,
							rsp.Asset,
							rsp.Pair.String(),
							rsp.ExecutionStrategy,
							rsp.Venue,
							float64(risk),
							float64(rsp.NotionalSize),
							rsp.Timestamp.AsTime(),
							rsp.SuccessfulOrders,
							rsp.Error,
						); err != nil {
							slog.Error(newCtx, "Failed to notify user of successful trade strategy: %v TradeParticipantId: %v", in.TradeStrategyId, rsp.TradeParticipantId)
						}

						// Push to pulse channel.
						if err := notifyPulseChannelUserTradeSuccess(newCtx, userID, in.TradeStrategyId, rsp.ExecutionStrategy, rsp.Venue, risk, rsp.SuccessfulOrders); err != nil {
							slog.Error(newCtx, err.Error())
						}
					}
				}

			case <-childCtx.Done():
				slog.Warn(newCtx, "Closing window for new trade participants for trade strategy: %v", in.TradeStrategyId)

				if err := notifyTradesChannelContextEnded(newCtx, in.TradeStrategyId); err != nil {
					slog.Error(newCtx, err.Error())
				}

				if err := notifyPulseChannelEnd(newCtx, in.TradeStrategyId, deadline); err != nil {
					slog.Error(newCtx, err.Error())
				}

				return
			}
		}
	}()

	return &satoshiproto.PollTradeStrategyParticipantsResponse{}, nil
}
