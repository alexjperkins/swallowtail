package handler

import (
	"context"
	"fmt"
	"strconv"
	"swallowtail/libraries/gerrors"
	discordproto "swallowtail/s.discord/proto"
)

const (
	maxCharacterPerMsg = 2000
)

// SendBatchMsgToChannel ...
func (s *DiscordService) SendBatchMsgToChannel(
	ctx context.Context, in *discordproto.SendBatchMsgToChannelRequest,
) (*discordproto.SendBatchMsgToChannelResponse, error) {
	switch {
	case in.Content == "":
		return nil, gerrors.BadParam("missing_param.content", nil)
	case in.ChannelId == "":
		return nil, gerrors.BadParam("missing_param.channel_id", nil)
	}

	errParams := map[string]string{
		"channel_id":      in.ChannelId,
		"idempotency_key": in.IdempotencyKey,
		"sender_id":       in.SenderId,
		"force":           strconv.FormatBool(in.Force),
	}

	if len(in.Content) < maxCharacterPerMsg {
		if _, err := (&discordproto.SendMsgToChannelRequest{
			Content:        in.Content,
			ChannelId:      in.ChannelId,
			IdempotencyKey: in.IdempotencyKey,
			SenderId:       in.SenderId,
			Force:          in.Force,
		}).Send(ctx).Response(); err != nil {
			return nil, gerrors.Augment(err, "failed_to_send_batch_msg_to_channel", errParams)
		}

		return &discordproto.SendBatchMsgToChannelResponse{}, nil
	}

	var msgs []string
	switch in.Separator {
	case "":
		var err error
		msgs, err = emptySeparatorHandler(in.Content)
		if err != nil {
			return nil, gerrors.Augment(err, "failed_to_send_batch_msg_to_channel", errParams)
		}
	default:
		var err error
		msgs, err = nonEmptySeparatorHandler(in.Content, in.Separator)
		if err != nil {
			return nil, gerrors.Augment(err, "failed_to_send_batch_msg_to_channel", errParams)
		}
	}

	for i, msg := range msgs {
		// We do this sequentially as to keep the order of msgs.
		if _, err := (&discordproto.SendMsgToChannelRequest{
			Content:        msg,
			ChannelId:      in.ChannelId,
			IdempotencyKey: fmt.Sprintf("%s-%d", in.IdempotencyKey, i),
			SenderId:       in.SenderId,
			Force:          in.Force,
		}).Send(ctx).Response(); err != nil {
			errParams["msg_num"] = strconv.Itoa(i)
			return nil, gerrors.Augment(err, "failed_to_send_batch_msg_to_channel", errParams)
		}
	}

	return &discordproto.SendBatchMsgToChannelResponse{}, nil
}
