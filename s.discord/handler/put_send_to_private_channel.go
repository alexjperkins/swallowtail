package handler

import (
	"swallowtail/s.discord/client"
	"swallowtail/s.discord/dao"
	"swallowtail/s.discord/domain"
	discordproto "swallowtail/s.discord/proto"
	"time"

	"github.com/monzo/terrors"
	"github.com/monzo/typhon"
)

func PUTSendToPrivateChannel(req typhon.Request) typhon.Response {
	body := discordproto.SendMsgToPrivateChannelRequest{}
	if err := req.Decode(body); err != nil {
		return typhon.Response{Error: err}
	}

	errParams := map[string]string{
		"idempotency_key": body.IdempotencyKey,
		"user_id":         body.UserId,
		"sender_id":       body.SenderId,
	}

	// First lets check if the idempotency key exists in persistent storage.
	_, exists, err := dao.Exists(req, body.IdempotencyKey)
	if err != nil {
		return typhon.Response{Error: terrors.Augment(err, "Failed to read existing; dao failed to read", errParams)}
	}
	switch {
	case exists && !body.Force:
		return req.Response(&discordproto.SendMsgPrivateToChannelResponse{})
	}

	// Send message via discord.
	if err := client.SendPrivateMessage(req, body.Content, body.UserId); err != nil {
		return typhon.Response{Error: terrors.Augment(err, "Failed to send message via discord.", errParams)}
	}

	// If the touch doesn't exist or the sender wants to force through an update; then we set via the dao.
	switch {
	case !exists:
		if _, err := (dao.Create(req, &domain.Touch{
			IdempotencyKey: body.IdempotencyKey,
			SenderID:       body.SenderId,
			Updated:        time.Now(),
		})); err != nil {
			// We do have the case whereby the write fails but we still send the message; this is preferable
			// to persisting the idempotency key, but failing to send.
			// We can take the hit of duplicate messages.
			return typhon.Response{Error: terrors.Augment(err, "Failed to create touch.", errParams)}
		}
	default:
		if _, err := (dao.Update(req, &domain.Touch{
			IdempotencyKey: body.IdempotencyKey,
			SenderID:       body.SenderId,
			Updated:        time.Now(),
		})); err != nil {
			// We have the same case as above here too.
			return typhon.Response{Error: terrors.Augment(err, "Failed to update touch.", errParams)}
		}
	}

	return req.Response(nil)
}
