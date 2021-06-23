package handler

import (
	discordproto "swallowtail/s.discord/proto"
)

// DiscordService implements the service for discord.
type DiscordService struct {
	discordproto.UnimplementedDiscordServer
}
