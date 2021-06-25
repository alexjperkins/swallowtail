package handler

import (
	"context"
	"swallowtail/s.discord/dao"
	"swallowtail/s.discord/domain"
	"time"

	"github.com/monzo/terrors"

	"swallowtail/s.discord/client"
	discordproto "swallowtail/s.discord/proto"
)

// PUTSendToChannel gPRC handler for sending messages to a given channel via discord.
func (s *DiscordService) PUTSendToChannel(
	ctx context.Context, in *discordproto.SendMsgToChannelRequest,
) (*discordproto.SendMsgToChannelResponse, error) {

	errParams := map[string]string{
		"idempotency_key": in.IdempotencyKey,
		"channel_id":      in.ChannelId,
		"sender_id":       in.SenderId,
	}

	// First lets check if the idempotency key exists in persistent storage.
	_, exists, err := dao.Exists(ctx, in.IdempotencyKey)
	if err != nil {
		return nil, terrors.Augment(err, "Failed to read existing; dao failed to read", errParams)
	}
	switch {
	case exists && !in.Force:
		return &discordproto.SendMsgToChannelResponse{}, nil
	}

	// Send message via discord.
	if err := client.Send(ctx, in.Content, in.ChannelId); err != nil {
		return nil, terrors.Augment(err, "Failed to send message via discord.", errParams)
	}

	// If the touch doesn't exist or the sender wants to force through an update; then we set via the dao.
	switch {
	case !exists:
		if _, err := (dao.Create(ctx, &domain.Touch{
			IdempotencyKey: in.IdempotencyKey,
			SenderID:       in.SenderId,
			Updated:        time.Now(),
		})); err != nil {
			// We do have the case whereby the write fails but we still send the message; this is preferable
			// to persisting the idempotency key, but failing to send.
			// We can take the hit of duplicate messages.
			return nil, terrors.Augment(err, "Failed to create touch.", errParams)
		}
	default:
		if _, err := (dao.Update(ctx, &domain.Touch{
			IdempotencyKey: in.IdempotencyKey,
			SenderID:       in.SenderId,
			Updated:        time.Now(),
		})); err != nil {
			// We have the same case as above here too.
			return nil, terrors.Augment(err, "Failed to update touch.", errParams)
		}
	}

	return &discordproto.SendMsgToChannelResponse{}, nil
}
