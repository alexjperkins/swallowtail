package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/monzo/slog"

	"swallowtail/libraries/gerrors"
	"swallowtail/s.binance/client"
	binanceproto "swallowtail/s.binance/proto"
	discordproto "swallowtail/s.discord/proto"
)

// GetStatus ...
func (s *BinanceService) GetStatus(
	ctx context.Context, in *binanceproto.GetStatusRequest,
) (*binanceproto.GetStatusResponse, error) {
	rsp, err := client.GetStatus(ctx)
	if err != nil {
		notifyExchangePulseChannelFailure(ctx)

		return nil, gerrors.Augment(err, "failed_to_get_status", nil)
	}

	// Convert to milliseconds.
	latencyInMilliseconds := int(rsp.ServerLatency / time.Millisecond)
	driftInMilliseconds := int(rsp.AssumedClockDrift / time.Millisecond)

	notifyExchangePulseChannelSuccess(ctx, rsp.ServerTime, latencyInMilliseconds, driftInMilliseconds)

	return &binanceproto.GetStatusResponse{
		ServerTime:        int64(rsp.ServerTime),
		RequestLatency:    int64(latencyInMilliseconds),
		AssumedClockDrift: int64(rsp.AssumedClockDrift),
	}, nil
}

func notifyExchangePulseChannelFailure(ctx context.Context) error {
	header := ":rotating_lights:    `CRONITOR: BINANCE CONNECTIVITY DOWN`    :rotating_lights:"
	content := `
SUCCESS: FALSE
`
	_, err := (&discordproto.SendMsgToChannelRequest{
		ChannelId:      discordproto.DiscordSatoshiExchangePulseChannel,
		SenderId:       "c.exchange",
		Content:        fmt.Sprintf("%s```%s```", header, content),
		IdempotencyKey: fmt.Sprintf("%s-%s", "binancestatus", time.Now().UTC().Truncate(time.Minute).String()),
	}).Send(ctx).Response()
	if err != nil {
		// Best Effort
		slog.Error(ctx, "Failed to notify pulse channel for failure: Error: %v", err)
	}

	return nil
}

func notifyExchangePulseChannelSuccess(ctx context.Context, serverTime, drift, latency int) error {
	header := ":taxi:    `CRONITOR: EXCHANGE PULSE`    :taxi:"
	content := `
EXCHANGE:     BINANCE
SUCCESS:      TRUE 
SERVERTIME:   %v
DRIFT (ms):   %v
LATENCY (ms): %v
`
	formattedContent := fmt.Sprintf(content, time.Unix(int64(serverTime), 0), drift, latency)

	_, err := (&discordproto.SendMsgToChannelRequest{
		ChannelId:      discordproto.DiscordSatoshiExchangePulseChannel,
		SenderId:       "c.exchange",
		Content:        fmt.Sprintf("%s```%s```", header, formattedContent),
		IdempotencyKey: fmt.Sprintf("%s-%s", "binancestatus", time.Now().UTC().Truncate(time.Minute).String()),
	}).Send(ctx).Response()
	if err != nil {
		// Best Effort
		slog.Error(ctx, "Failed to notify pulse channel for success: Error: %v", err)
	}

	return nil
}
