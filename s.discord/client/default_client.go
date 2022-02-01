package client

import (
	"context"
	"fmt"

	"swallowtail/libraries/gerrors"
	"swallowtail/s.discord/domain"

	"github.com/bwmarrin/discordgo"
	"github.com/monzo/slog"
	"github.com/monzo/terrors"

	discordproto "swallowtail/s.discord/proto"
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

	// Open websocket session.
	if err = s.Open(); err != nil {
		panic(err)
	}

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

func (d *discordClient) Send(ctx context.Context, message, channelID string) (*discordgo.Message, error) {
	var cID = channelID
	if !d.isActive {
		cID = discordTestingChannel

	}

	msg, err := d.session.ChannelMessageSend(cID, message)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_send_message", nil)
	}

	return msg, nil
}

func (d *discordClient) SendPrivateMessage(ctx context.Context, message, userID string) error {
	if !d.isActive {
		// If not active; we simply send to the testing channel.
		_, err := d.Send(ctx, discordTestingChannel, message)
		return err
	}

	ch, err := d.session.UserChannelCreate(userID)
	if err != nil {
		return gerrors.Augment(err, "failed_to_create_private_channel", map[string]string{
			"discord_user_id": userID,
		})
	}

	_, err = d.Send(ctx, message, ch.ID)
	return err
}

func (d *discordClient) ReadRoles(ctx context.Context, userID string) ([]*domain.Role, error) {
	m, err := d.session.GuildMember(discordproto.DiscordSatoshiGuildID, userID)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_read_roles.failed_to_fetch_member", map[string]string{
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

func (d *discordClient) SetRoles(ctx context.Context, userID string, roles []*domain.Role) error {
	roleIDs := []string{}
	for _, role := range roles {
		roleIDs = append(roleIDs, role.ID)
	}

	if err := d.session.GuildMemberEdit(discordproto.DiscordSatoshiGuildID, userID, roleIDs); err != nil {
		return gerrors.Augment(err, "failed_to_set_roles", map[string]string{
			"guild_id": discordproto.DiscordSatoshiGuildID,
		})
	}

	return nil
}

func (d *discordClient) ReadMessageReactions(ctx context.Context, messageID, channelID string) (map[EmojiID][]string, error) {
	m, err := d.session.ChannelMessage(channelID, messageID)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_read_message_reactions.message", nil)
	}

	reactions := map[EmojiID][]string{}
	for _, reaction := range m.Reactions {
		users, err := d.session.MessageReactions(channelID, messageID, reaction.Emoji.Name, 100, "", "")
		if err != nil {
			return nil, gerrors.Augment(err, "failed_to_read_message_reactions.reactions", nil)
		}

		userIDs := []string{}
		for _, user := range users {
			userIDs = append(userIDs, user.ID)
		}

		reactions[EmojiID(reaction.Emoji.APIName())] = userIDs
	}

	return reactions, nil
}

func (d *discordClient) AddHandler(handler func(s *discordgo.Session, m *discordgo.MessageCreate)) {
	d.session.AddHandler(handler)
}

func (d *discordClient) AddHandlerGuildMemberAdd(handler func(s *discordgo.Session, u *discordgo.GuildMemberAdd)) {
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
