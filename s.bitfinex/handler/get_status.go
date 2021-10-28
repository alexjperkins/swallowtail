package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/monzo/slog"

	"swallowtail/libraries/gerrors"
	"swallowtail/s.bitfinex/client"
	"swallowtail/s.bitfinex/dto"
	"swallowtail/s.bitfinex/marshaling"
	bitfinexproto "swallowtail/s.bitfinex/proto"
	discordproto "swallowtail/s.discord/proto"
)

// GetStatus ...
func (s *BitfinexService) GetStatus(
	ctx context.Context, in *bitfinexproto.GetBitfinexStatusRequest,
) (*bitfinexproto.GetBitfinexStatusResponse, error) {
	// Get Bitfinex status.
	rsp, err := client.GetStatus(ctx, &dto.GetStatusRequest{})
	if err != nil {
		// Notify of failure.
		notifyPulseChannelOnFailure(ctx)
		return nil, gerrors.Augment(err, "failed_to_get_status", nil)
	}

	// Notify of success.
	notifyPulseChannelOnSuccess(ctx, rsp)

	// Marshal from DTO to proto.
	proto := marshaling.GetStatusDTOToProto(rsp)

	return proto, nil
}

func notifyPulseChannelOnFailure(ctx context.Context) error {
	header := ":rotating_lights:    `CRONITOR: BITFINEX CONNECTIVITY DOWN`    :rotating_lights:"
	content := `
SUCCESS: FALSE
`
	_, err := (&discordproto.SendMsgToChannelRequest{
		ChannelId:      discordproto.DiscordSatoshiExchangePulseChannel,
		SenderId:       "c.exchange",
		Content:        fmt.Sprintf("%s```%s```", header, content),
		IdempotencyKey: fmt.Sprintf("%s-%s", "bitfinexstatus", time.Now().UTC().Truncate(time.Minute).String()),
	}).Send(ctx).Response()
	if err != nil {
		// Best Effort
		slog.Error(ctx, "Failed to notify pulse channel for failure: Error: %v", err)
	}

	return nil
}

func notifyPulseChannelOnSuccess(ctx context.Context, serverRsp *dto.GetStatusResponse) error {
	header := ":rabbit:    `CRONITOR: EXCHANGE PULSE`    :rabbit2:"
	content := `
EXCHANGE:     BITFINEX
SUCCESS:      TRUE 
LATENCY (ms): %v
`
	formattedContent := fmt.Sprintf(content, serverRsp.ServerLatency)

	_, err := (&discordproto.SendMsgToChannelRequest{
		ChannelId:      discordproto.DiscordSatoshiExchangePulseChannel,
		SenderId:       "c.exchange",
		Content:        fmt.Sprintf("%s```%s```", header, formattedContent),
		IdempotencyKey: fmt.Sprintf("%s-%s", "bitfinexstatus", time.Now().UTC().Truncate(time.Minute).String()),
	}).Send(ctx).Response()
	if err != nil {
		// Best Effort
		slog.Error(ctx, "Failed to notify pulse channel for success: Error: %v", err)
	}

	return nil
}
