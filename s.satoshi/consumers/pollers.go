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

func startTradeParticipantsPoller(ctx context.Context, messageID, tradeID string) error {
	if _, err := (&satoshiproto.PollTradeParticipantsRequest{
		ActorId:          satoshiproto.SatoshiActorSatoshiSystem,
		TradeId:          tradeID,
		ChannelId:        discordproto.DiscordSatoshiModTradesChannel,
		MessageId:        messageID,
		TimeoutInMinutes: int64(tradeParticipantTimeout),
	}).Send(ctx).Response(); err != nil {
		return gerrors.Augment(err, "failed_to_start_participants_poller", map[string]string{
			"message_id": messageID,
			"trade_id":   tradeID,
		})
	}

	return nil
}
