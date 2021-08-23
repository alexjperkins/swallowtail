package pager

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"swallowtail/libraries/gerrors"
	"swallowtail/libraries/util"
	accountproto "swallowtail/s.account/proto"
	discordproto "swallowtail/s.discord/proto"
)

type discordPager struct{}

func init() {
	register(accountproto.PagerType_DISCORD.String(), &discordPager{})
}

func (d *discordPager) Page(ctx context.Context, userID, msg string) error {
	hashedContent := util.Sha256Hash(msg)
	now := time.Now()

	if _, err := (&discordproto.SendMsgToPrivateChannelRequest{
		UserId:   userID,
		Content:  msg,
		SenderId: "system:s.account:pager",
		// Idempotent on channel, message & the hour of the day.
		IdempotencyKey: fmt.Sprintf("%s-%s-%s", userID, hashedContent, strconv.Itoa(now.Hour())),
	}).Send(ctx).Response(); err != nil {
		return gerrors.Augment(err, "failed_to_page_user", map[string]string{
			"user_id": userID,
		})
	}

	return nil
}
