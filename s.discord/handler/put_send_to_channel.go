package handler

import (
	"context"
	"time"

	"swallowtail/libraries/gerrors"
	"swallowtail/s.discord/client"
	"swallowtail/s.discord/dao"
	"swallowtail/s.discord/domain"
	discordproto "swallowtail/s.discord/proto"
)

// SendMsgToChannel gRPC handler for sending messages to a given channel via discord.
func (s *DiscordService) SendMsgToChannel(
	ctx context.Context, in *discordproto.SendMsgToChannelRequest,
) (*discordproto.SendMsgToChannelResponse, error) {
	switch {
	case in.ChannelId == "":
		return nil, gerrors.BadParam("missing_param.channel_id", nil)
	case in.Content == "":
		return nil, gerrors.BadParam("missing_param.content", nil)
	}

	errParams := map[string]string{
		"idempotency_key": in.IdempotencyKey,
		"channel_id":      in.ChannelId,
		"sender_id":       in.SenderId,
	}

	var exists bool
	if in.IdempotencyKey != "" {
		// First lets check if the idempotency key exists in persistent storage.
		_, doesExist, err := dao.Exists(ctx, in.IdempotencyKey)
		if err != nil {
			return nil, gerrors.Augment(err, "failed_to_send_msg_to_channel.failed_to_read_idempotency_key", errParams)
		}
		exists = doesExist
	}

	switch {
	case exists && !in.Force:
		return &discordproto.SendMsgToChannelResponse{}, nil
	}

	// Send message via discord.
	if err := client.Send(ctx, in.Content, in.ChannelId); err != nil {
		return nil, gerrors.Augment(err, "failed_to_send_msg_channel.client_failure", errParams)
	}

	// If the touch doesn't exist or the sender wants to force through an update; then we set via the dao.
	switch {
	case !exists && in.IdempotencyKey != "":
		if _, err := (dao.Create(ctx, &domain.Touch{
			IdempotencyKey: in.IdempotencyKey,
			SenderID:       in.SenderId,
			Updated:        time.Now(),
		})); err != nil {
			// We do have the case whereby the write fails but we still send the message; this is preferable
			// to persisting the idempotency key, but failing to send.
			// We can take the hit of duplicate messages.
			return nil, gerrors.Augment(err, "failed_to_touch_discord_message.", errParams)
		}
	default:
		if _, err := (dao.Update(ctx, &domain.Touch{
			IdempotencyKey: in.IdempotencyKey,
			SenderID:       in.SenderId,
			Updated:        time.Now(),
		})); err != nil {
			// We have the same case as above here too.
			return nil, gerrors.Augment(err, "failed_to_update_touch_discord_message.", errParams)
		}
	}

	return &discordproto.SendMsgToChannelResponse{}, nil
}
