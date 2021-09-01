package client

import (
	"context"
	"fmt"

	"swallowtail/libraries/gerrors"
	"swallowtail/s.discord/domain"
	discordproto "swallowtail/s.discord/proto"

	"github.com/bwmarrin/discordgo"
	"github.com/monzo/slog"
	"github.com/monzo/terrors"
)

// New creates a new discord client
func New(name, token string, isBot bool) DiscordClient {
	t := formatToken(token, isBot)
	s, err := discordgo.New(t)
	if err != nil {
		panic(terrors.Augment(err, "Failed to create discord client", map[string]string{
			"discord_token": t,
			"name":          name,
		}))
	}

	s.LogLevel = 0
	s.Debug = true

	// Open websocket session.
	err = s.Open()
	if err != nil {
		panic(err)
	}

	slog.Info(nil, "USER-AGENT: %s", s.UserAgent)
	if !isActiveFlag {
		slog.Warn(context.TODO(), "Discord client set to TESTING MODE.")
	}

	slog.Info(context.TODO(), "Created discord bot: %s, token: %s", name, t)
	return &discordClient{
		session:  s,
		isBot:    isBot,
		isActive: isActiveFlag,
	}
}

type discordClient struct {
	session  *discordgo.Session
	isBot    bool
	isActive bool
}

func (d *discordClient) Send(ctx context.Context, message, channelID string) error {
	var (
		cID = channelID
	)
	if !d.isActive {
		cID = discordTestingChannel

	}
	msg, err := d.session.ChannelMessageSend(cID, message)
	if err != nil {
		return err
	}
	slog.Info(ctx, "Message Posted to discord: %v", msg)
	return nil
}

func (d *discordClient) SendPrivateMessage(ctx context.Context, message, userID string) error {
	if !d.isActive {
		// If not active; we simply send to the testing channel.
		return d.Send(ctx, discordTestingChannel, message)
	}

	slog.Info(ctx, "user-id: %s", userID)

	ch, err := d.session.UserChannelCreate(userID)
	if err != nil {
		return terrors.Augment(err, "Failed to create private channel", map[string]string{
			"discord_user_id": userID,
		})
	}

	return d.Send(ctx, message, ch.ID)
}

func (d *discordClient) ReadRoles(ctx context.Context, userID string) ([]*domain.Role, error) {
	m, err := d.session.GuildMember(discordproto.DiscordSatoshiGuildID, userID)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_read_roles", map[string]string{
			"guild_id": discordproto.DiscordSatoshiGuildID,
		})
	}

	roles := []*domain.Role{}
	for _, r := range m.Roles {
		name, ok := discordproto.ConvertRoleIDToName(r)
		if !ok {
			slog.Warn(ctx, "Invalid role ID: %s", r)
			continue
		}

		roles = append(roles, &domain.Role{
			ID:   r,
			Name: name,
		})
	}

	return roles, nil
}

func (d *discordClient) AddHandler(handler func(s *discordgo.Session, m *discordgo.MessageCreate)) {
	slog.Info(nil, "Adding handler")
	d.session.AddHandler(handler)
}

func (d *discordClient) Close() {
	d.session.Close()
}

func (d *discordClient) Ping(ctx context.Context) error {
	// TODO: best way to ping the discord client?
	return nil
}

func formatToken(token string, isBot bool) string {
	if !isBot {
		return token
	}
	return fmt.Sprintf("Bot %s", token)
}
