package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/monzo/slog"

	discordproto "swallowtail/s.discord/proto"
	"swallowtail/s.ftx/client"
	ftxproto "swallowtail/s.ftx/proto"
)

// GetFTXStatus ...
func (s *FTXService) GetFTXStatus(
	ctx context.Context, _ *ftxproto.GetFTXStatusRequest,
) (*ftxproto.GetFTXStatusResponse, error) {
	then := time.Now().UTC()

	rsp, err := client.GetStatus(ctx, nil)
	if err != nil {
		notifyExchangePulseChannelFailure(ctx)
		return nil, err
	}

	requestLatency := int(time.Since(then) / time.Millisecond)

	var p50Latency float64
	if len(rsp.Result) > 0 {
		p50Latency = rsp.Result[0].P50Latency
	}

	notifyExchangePulseChannelSuccess(ctx, requestLatency, int(p50Latency*1000))

	return &ftxproto.GetFTXStatusResponse{}, nil
}

func notifyExchangePulseChannelFailure(ctx context.Context) error {
	header := ":rotating_lights:    `CRONITOR: FTX CONNECTIVITY DOWN`    :rotating_lights:"
	content := `
SUCCESS: FALSE
`
	_, err := (&discordproto.SendMsgToChannelRequest{
		ChannelId:      discordproto.DiscordSatoshiExchangePulseChannel,
		SenderId:       "c.exchange",
		Content:        fmt.Sprintf("%s```%s```", header, content),
		IdempotencyKey: fmt.Sprintf("%s-%s", "ftxstatus", time.Now().UTC().Truncate(time.Minute).String()),
	}).Send(ctx).Response()
	if err != nil {
		// Best Effort
		slog.Error(ctx, "Failed to notify pulse channel for failure: Error: %v", err)
	}

	return nil
}

func notifyExchangePulseChannelSuccess(ctx context.Context, requestLatency, p50Latency int) error {
	header := ":taxi:    `CRONITOR: EXCHANGE PULSE`    :taxi:"
	content := `
EXCHANGE:     FTX 
SUCCESS:      TRUE 
LATENCY (ms): %v
P50     (ms): %v

	return nil
`
	formattedContent := fmt.Sprintf(content, requestLatency, p50Latency)

	_, err := (&discordproto.SendMsgToChannelRequest{
		ChannelId:      discordproto.DiscordSatoshiExchangePulseChannel,
		SenderId:       "c.exchange",
		Content:        fmt.Sprintf("%s```%s```", header, formattedContent),
		IdempotencyKey: fmt.Sprintf("%s-%s", "ftxstatus", time.Now().UTC().Truncate(time.Minute).String()),
	}).Send(ctx).Response()
	if err != nil {
		// Best Effort
		slog.Error(ctx, "Failed to notify pulse channel for success: Error: %v", err)
	}

	return nil
}
