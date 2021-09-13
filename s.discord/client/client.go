package client

import (
	"context"
	"fmt"

	"swallowtail/libraries/util"
	"swallowtail/s.discord/domain"

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

type EmojiID string

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
	// Send
	Send(ctx context.Context, message, channelID string) (*discordgo.Message, error)
	// SendPrivateMessage
	SendPrivateMessage(ctx context.Context, message, userID string) error
	// AddHandler
	AddHandler(handler func(s *discordgo.Session, m *discordgo.MessageCreate))
	// ReadRoles
	ReadRoles(ctx context.Context, userID string) ([]*domain.Role, error)
	// SetRoles set the users roles to the roles passed. It replaces all the roles the user currently has.
	SetRoles(ctx context.Context, userID string, roles []*domain.Role) error
	// ReadMessageReactions returns a map of reactions by emoji id and the list of users who made that reaction.
	ReadMessageReactions(ctx context.Context, messageID, channelID string) (map[EmojiID][]string, error)

	// TODO: Remove below
	Close()
	Ping(ctx context.Context) error
}

// Send sends a message to a given channel`channel_id` via discord.
func Send(ctx context.Context, message, channelID string) (*discordgo.Message, error) {
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

// ReadRoles ...
func ReadRoles(ctx context.Context, userID string) ([]*domain.Role, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Read user roles")
	defer span.Finish()
	return client.ReadRoles(ctx, userID)
}

// SetRoles ...
func SetRoles(ctx context.Context, userID string, roles []*domain.Role) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Set user roles")
	defer span.Finish()
	return client.SetRoles(ctx, userID, roles)
}

// ReadMessage ...
func ReadMessageReactions(ctx context.Context, messageID, channelID string) (map[EmojiID][]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Read message")
	defer span.Finish()
	return client.ReadMessageReactions(ctx, messageID, channelID)
}
