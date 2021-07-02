package pager

import (
	"context"
	"fmt"
	"strconv"
	"time"

	accountproto "swallowtail/s.account/proto"
	discordproto "swallowtail/s.discord/proto"

	"github.com/monzo/terrors"

	"google.golang.org/grpc"
)

type discordPager struct{}

func init() {
	register(accountproto.PagerType_DISCORD.String(), &discordPager{})
}

func (d *discordPager) Page(ctx context.Context, channelID, msg string) error {
	now := time.Now()

	conn, err := grpc.DialContext(ctx, "s_discord")
	if err != nil {
		return terrors.Augment(err, "Failed to reach s.discord via rpc", nil)
	}
	defer conn.Close()

	client := discordproto.NewDiscordClient(conn)
	client.SendMsgToChannel(ctx, &discordproto.SendMsgToChannelRequest{
		ChannelId: channelID,
		Content:   msg,
		SenderId:  "system:s.account:pager",
		// Idempotent on channel, message & the hour of the day.
		IdempotencyKey: fmt.Sprintf("%s-%s-%s", channelID, hash(msg), strconv.Itoa(now.Hour())),
	})
	return nil
}
