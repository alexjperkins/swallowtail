package parser

import (
	discordproto "swallowtail/s.discord/proto"
)

func init() {
	register(discordproto.DiscordSatoshiInternalCallsChannel, []TradeParser{
		&DCAParser{},
		&DMAParser{},
	})
}
