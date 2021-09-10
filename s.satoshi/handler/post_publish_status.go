package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/monzo/slog"

	"swallowtail/libraries/gerrors"
	discordproto "swallowtail/s.discord/proto"
	satoshiproto "swallowtail/s.satoshi/proto"
	"swallowtail/s.satoshi/satoshi"
)

var (
	version = satoshi.Version
)

// PublishStatus ...
func (s *SatoshiService) PublishStatus(
	ctx context.Context, in *satoshiproto.PublishStatusRequest,
) (*satoshiproto.PublishStatusResponse, error) {
	if err := notifyPulseChannelSuccess(ctx, version); err != nil {
		return nil, gerrors.Augment(err, "failed_to_publish_status.notify_pulse_channel", nil)
	}

	return &satoshiproto.PublishStatusResponse{
		Alive:   true,
		Version: version,
	}, nil
}

func notifyPulseChannelSuccess(ctx context.Context, version string) error {
	header := "<:satoshi:886008008491024445>    `CRONITOR: SATOSHI PULSE`    <:satoshi:886008008491024445>"
	content := `
ALIVE:   TRUE
VERSION: %s 
`
	formattedContent := fmt.Sprintf(content, version)

	_, err := (&discordproto.SendMsgToChannelRequest{
		ChannelId:      discordproto.DiscordSatoshiSatoshiPulseChannel,
		SenderId:       "c.satoshi",
		Content:        fmt.Sprintf("%s```%s```", header, formattedContent),
		IdempotencyKey: fmt.Sprintf("%s-%s", "satoshistatus", time.Now().UTC().Truncate(time.Minute).String()),
	}).Send(ctx).Response()
	if err != nil {
		// Best Effort
		slog.Error(ctx, "Failed to notify pulse channel for success: Error: %v", err)
	}

	return nil
}
