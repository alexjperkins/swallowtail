package consumers

import (
	"context"

	"swallowtail/libraries/gerrors"
	discordproto "swallowtail/s.discord/proto"
	satoshiproto "swallowtail/s.satoshi/proto"
)

var (
	tradeParticipantTimeout = 15
)

func startTradeParticipantsPoller(ctx context.Context, messageID, tradeStrategyID string) error {
	if _, err := (&satoshiproto.PollTradeStrategyParticipantsRequest{
		ActorId:          satoshiproto.SatoshiActorSatoshiSystem,
		TradeStrategyId:  tradeStrategyID,
		ChannelId:        discordproto.DiscordSatoshiModTradesChannel,
		MessageId:        messageID,
		TimeoutInMinutes: int64(tradeParticipantTimeout),
	}).Send(ctx).Response(); err != nil {
		return gerrors.Augment(err, "failed_to_start_participants_poller", map[string]string{
			"message_id": messageID,
			"trade_id":   tradeStrategyID,
		})
	}

	return nil
}
