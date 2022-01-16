package parser

import (
	discordproto "swallowtail/s.discord/proto"
)

func init() {
	register(discordproto.DiscordMoonSwingGroupChannel, []TradeParser{
		&DCAParser{},
		&DMAParser{},
	})
}
