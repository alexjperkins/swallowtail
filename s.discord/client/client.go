package client

import (
	"context"
	"fmt"

	"swallowtail/libraries/util"

	"github.com/bwmarrin/discordgo"
	"github.com/monzo/slog"
	"github.com/monzo/terrors"
	"github.com/opentracing/opentracing-go"
)

var (
	// TODO: change implementation to use own defined mocks
	DiscordClientID = "discord-client-id"

	isActiveFlag          bool
	discordTestingChannel = "817513133274824715"

	client      DiscordClient
	clientToken string
)

func init() {
	clientToken = fmt.Sprintf("%v", util.EnvGetOrDefault("SATOSHI_DISCORD_API_TOKEN", ""))
	v := util.EnvGetOrDefault("DISCORD_TESTING_OVERRIDE", "0")
	if v != "1" {
		isActiveFlag = true
	}
}

func Init(ctx context.Context) error {
	// Won't work yet until we migrate to RPCs.
	c := New(DiscordClientID, clientToken, true)

	if err := c.Ping(ctx); err != nil {
		return terrors.Augment(err, "Failed to establish connection with discord client", nil)
	}

	slog.Info(ctx, "Discord client initialized", nil)

	client = c
	return nil
}

// DiscordClient ...
type DiscordClient interface {
	Send(ctx context.Context, message, channelID string) error
	SendPrivateMessage(ctx context.Context, message, userID string) error
	AddHandler(handler func(s *discordgo.Session, m *discordgo.MessageCreate))
	Close()
	Ping(ctx context.Context) error
}

// Send sends a message to a given channel`channel_id` via discord.
func Send(ctx context.Context, message, channelID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Send discord message via channel")
	defer span.Finish()
	return client.Send(ctx, message, channelID)
}

// Send sends a private message to a given user `user_id` via discord.
func SendPrivateMessage(ctx context.Context, message, userID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Send discord message via private channel")
	defer span.Finish()
	return client.SendPrivateMessage(ctx, message, userID)
}
