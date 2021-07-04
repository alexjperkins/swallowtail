package pager

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/monzo/terrors"
	"google.golang.org/grpc"

	"swallowtail/libraries/util"
	accountproto "swallowtail/s.account/proto"
	discordproto "swallowtail/s.discord/proto"
)

type discordPager struct{}

func init() {
	register(accountproto.PagerType_DISCORD.String(), &discordPager{})
}

func (d *discordPager) Page(ctx context.Context, userID, msg string) error {
	now := time.Now()

	conn, err := grpc.DialContext(ctx, "swallowtail-s-discord:8000", grpc.WithInsecure())
	if err != nil {
		return terrors.Augment(err, "Failed to reach s.discord via rpc", nil)
	}
	defer conn.Close()

	hashedContent, err := util.Sha256Hash(msg)
	if err != nil {
		return terrors.Augment(err, "Failed to page account; error hashing message content", map[string]string{
			"user_id": userID,
		})
	}

	client := discordproto.NewDiscordClient(conn)
	if _, err = (client.SendMsgToPrivateChannel(ctx, &discordproto.SendMsgToPrivateChannelRequest{
		UserId:   userID,
		Content:  msg,
		SenderId: "system:s.account:pager",
		// Idempotent on channel, message & the hour of the day.
		IdempotencyKey: fmt.Sprintf("%s-%s-%s", userID, hashedContent, strconv.Itoa(now.Hour())),
	})); err != nil {
		return terrors.Augment(err, "Failed to send msg to discord channel", map[string]string{
			"user_id": userID,
		})
	}

	return nil
}
