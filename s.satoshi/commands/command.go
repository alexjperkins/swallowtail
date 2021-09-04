package commands

import (
	"context"
	"fmt"
	"strings"
	"swallowtail/libraries/gerrors"
	"swallowtail/libraries/util"
	discordproto "swallowtail/s.discord/proto"

	"github.com/bwmarrin/discordgo"
	"github.com/monzo/slog"
)

const (
	commandIdentifier = "!"
)

// Command ...
type Command struct {
	ID string
	// IsPrivate dictates if the command must be run in a private channel.
	IsPrivate     bool
	IsFuturesOnly bool
	IsAdminOnly   bool
	// Non-inclusive of the prefix.
	MinimumNumberOfArgs int
	Usage               string
	Guide               string
	FailureMsg          string
	Handler             CommandHandler
	SubCommands         map[string]*Command
}

// CommandHandler ...
type CommandHandler func(ctx context.Context, tokens []string, s *discordgo.Session, m *discordgo.MessageCreate) error

// Exec executes the given command; recursing down the command tree if a subcommand is detected.
func (c *Command) Exec(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Check Prefix
	prefix := fmt.Sprintf("%s%s", commandIdentifier, c.ID)
	if !strings.HasPrefix(m.Content, prefix) {
		return
	}

	ctx := context.Background()
	tokens := strings.Fields(normalizeContent(m.Content))

	slog.Trace(ctx, "Received command: %s with args: %v", c.ID, tokens[1:])

	if err := c.exec(ctx, tokens[1:], s, m); err != nil {
		slog.Info(ctx, "Parent command %s, failed with error: %v", c.ID, err)
	}
}

func (c *Command) exec(ctx context.Context, tokens []string, s *discordgo.Session, m *discordgo.MessageCreate) error {
	// Check Privacy.
	if c.IsPrivate && m.GuildID != "" {
		_, err := s.ChannelMessageSend(m.ChannelID, formatNonPublicMsg(m.Author.ID))
		return gerrors.Augment(err, "failed_to_page_user.private", nil)
	}

	// Get members from the guild; this is for when the user messages a command via a private channel.
	// Add extra safety in the fact that we reject the user if they try and message satoshi directly,
	// when they're not in the guild.
	membersRoles, err := getMembersRolesFromGuild(s, m.Author.ID)
	if err != nil {
		return gerrors.Augment(err, "failed_exec_command", nil)
	}

	// Check if they are indeed an admin member if the command requires so.
	if c.IsAdminOnly && !isAdmin(membersRoles) {
		_, err := s.ChannelMessageSend(m.ChannelID, formatNonAdminMsg(m.Author.ID))
		return gerrors.Augment(err, "failed_to_page_user.non_admin", nil)
	}

	// Check if they are a indeed a futures member if the command requires so.
	if c.IsFuturesOnly && !isFuturesMember(membersRoles) {
		_, err := s.ChannelMessageSend(m.ChannelID, formatNonFuturesMsg(m.Author.ID))
		return gerrors.Augment(err, "failed_to_page_user.non_futures_member", nil)
	}

	// Check Usage.
	if len(tokens) > 0 && strings.ToLower(tokens[0]) == "help" {
		_, err := s.ChannelMessageSend(m.ChannelID, formatUsageMsg(m.Author.ID, c.Usage, c.Guide))
		return gerrors.Augment(err, "failed_to_page_user.help", nil)
	}

	// Check we have at least the correct number of arguments to execute the command.
	if len(tokens) < c.MinimumNumberOfArgs {
		_, err := s.ChannelMessageSend(m.ChannelID, formatUsageMsg(m.Author.ID, c.Usage, c.Guide))
		return gerrors.Augment(err, "failed_to_page_user.bad_params", nil)
	}

	// If we have no args; then we must not have any subcommand; so let's try the parent command default.
	if len(tokens) == 0 {
		if err := c.Handler(ctx, tokens, s, m); err != nil {
			_, err := s.ChannelMessageSend(m.ChannelID, formatFailureMsg(m.Author.ID, c.FailureMsg, err))
			return gerrors.Augment(err, "failed_to_page_user.command_failure_no_tokens", nil)
		}

		return nil
	}

	// IF we don't have a subcommand that matches the "second" token; then we can
	// try to run the parent command instead.
	subCommand, ok := c.SubCommands[tokens[0]]
	if !ok {
		if err := c.Handler(ctx, tokens, s, m); err != nil {
			_, err := s.ChannelMessageSend(m.ChannelID, formatFailureMsg(m.Author.ID, c.FailureMsg, err))
			return gerrors.Augment(err, "failed_to_page_user.command_failure", nil)
		}

		return nil
	}

	// We have a subcommand; so let's execute it.
	if err := subCommand.exec(ctx, tokens[1:], s, m); err != nil {
		return err
	}

	return nil
}

func formatUsageMsg(userID, usage string, guide string) string {
	var formattedGuide string
	if guide != "" {
		formattedGuide = fmt.Sprintf("Guide: `%s`", guide)
	}

	return fmt.Sprintf(":wave: <@%s>\n%s\n%s", userID, formattedGuide, util.WrapAsCodeBlock(usage))
}

func formatNonAdminMsg(userID string) string {
	return fmt.Sprintf(":wave: <@%s>, apologies! But this command can only be run by admins :disappointed:", userID)
}

func formatNonFuturesMsg(userID string) string {
	return fmt.Sprintf(":wave: <@%s>, apologies! But this command can only be run by futures members :grimacing: Ping @ajperkins if you want to know how to become one", userID)
}

func formatNonPublicMsg(userID string) string {
	return fmt.Sprintf(":wave: <@%s>, Please DM satoshi this command instead, the response may contain sensitive information. Thanks", userID)
}

func formatFailureMsg(userID, failureMsg string, err error) string {
	var errMsg = err.Error()
	switch {
	case gerrors.Is(err, gerrors.ErrUnimplemented):
		errMsg = "Command unimplemented"
	}

	return fmt.Sprintf(
		":disappointed: Sorry <@%s>, I failed to execute that command.\n%s\n Error: %s\n.",
		userID, failureMsg, errMsg,
	)
}

func isFuturesMember(roles []string) bool {
	for _, role := range roles {
		if role == discordproto.DiscordSatoshiFuturesRoleID {
			return true
		}
	}

	return false
}

func isAdmin(roles []string) bool {
	for _, role := range roles {
		if role == discordproto.DiscordSatoshiAdminRoleID {
			return true
		}
	}

	return false
}

// Placeholder
func normalizeContent(content string) string {
	return content
}

func getMembersRolesFromGuild(session *discordgo.Session, userID string) ([]string, error) {
	m, err := session.GuildMember(discordproto.DiscordSatoshiGuildID, userID)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_read_members_guild_roles", nil)
	}

	return m.Roles, nil
}
