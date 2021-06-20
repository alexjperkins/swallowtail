package pager

import (
	"context"

	accountproto "swallowtail/s.account/proto"
)

type discordPager struct{}

func init() {
	register(accountproto.AccountPagerTypeDiscord, &discordPager{})
}

func (d *discordPager) Page(ctx context.Context, identifier, msg string) error {
	return nil
}
