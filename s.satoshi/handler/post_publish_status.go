package handler

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/monzo/slog"

	"swallowtail/libraries/gerrors"
	discordproto "swallowtail/s.discord/proto"
	"swallowtail/s.satoshi/consumers"
	satoshiproto "swallowtail/s.satoshi/proto"
)

// PublishStatus ...
func (s *SatoshiService) PublishStatus(
	ctx context.Context, in *satoshiproto.PublishStatusRequest,
) (*satoshiproto.PublishStatusResponse, error) {
	cs := consumers.Registry()

	if err := notifyPulseChannelSuccess(ctx, satoshiproto.SatoshiVersion, cs); err != nil {
		return nil, gerrors.Augment(err, "failed_to_publish_status.notify_pulse_channel", nil)
	}

	return &satoshiproto.PublishStatusResponse{
		Alive:   true,
		Version: satoshiproto.SatoshiVersion,
	}, nil
}

func notifyPulseChannelSuccess(ctx context.Context, version string, consumers map[string]consumers.Consumer) error {
	header := "<:satoshi:886008008491024445>    `CRONITOR: SATOSHI PULSE`    <:satoshi:886008008491024445>"
	content := `
ALIVE:   TRUE
VERSION: %s 
`
	formattedContent := fmt.Sprintf(content, version)

	var sb strings.Builder
	for id, c := range consumers {
		sb.WriteString(fmt.Sprintf("\n%s: %v", strings.ToUpper(id), c.IsActive()))
	}

	withConsumers := fmt.Sprintf("%s%s", formattedContent, sb.String())

	_, err := (&discordproto.SendMsgToChannelRequest{
		ChannelId:      discordproto.DiscordSatoshiSatoshiPulseChannel,
		SenderId:       "c.satoshi",
		Content:        fmt.Sprintf("%s```%s```", header, withConsumers),
		IdempotencyKey: fmt.Sprintf("%s-%s", "satoshistatus", time.Now().UTC().Truncate(time.Minute).String()),
	}).Send(ctx).Response()
	if err != nil {
		// Best Effort
		slog.Error(ctx, "Failed to notify pulse channel for success: Error: %v", err)
	}

	return nil
}
