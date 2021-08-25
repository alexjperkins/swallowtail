package commands

import (
	"context"
	"fmt"
	"strings"
	"swallowtail/libraries/gerrors"
	"swallowtail/libraries/util"

	"github.com/bwmarrin/discordgo"
	"github.com/monzo/slog"
)

const (
	commandIdentifier = "!"
)

// Command ...
type Command struct {
	ID      string
	Private bool
	// Non-inclusive of the prefix.
	MinimumNumberOfArgs int
	Usage               string
	Helps               string
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
	if c.Private && m.GuildID != "" {
		_, err := s.ChannelMessageSend(m.ChannelID, formatNonPublicMsg(m.Author.ID))
		return gerrors.Augment(err, "failed_to_page_user.private", nil)
	}

	// Check Usage.
	if len(tokens) > 0 && strings.ToLower(tokens[0]) == "help" {
		_, err := s.ChannelMessageSend(m.ChannelID, formatUsageMsg(m.Author.ID, c.Usage))
		return gerrors.Augment(err, "failed_to_page_user.help", nil)
	}

	// Check we have at least the correct number of arguments to execute the command.
	if len(tokens) < c.MinimumNumberOfArgs {
		_, err := s.ChannelMessageSend(m.ChannelID, formatUsageMsg(m.Author.ID, c.Usage))
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

func formatUsageMsg(userID, usage string) string {
	return fmt.Sprintf(":wave: <@%s> %s", userID, util.WrapAsCodeBlock(usage))
}

func formatNonPublicMsg(userID string) string {
	return fmt.Sprintf(":wave: <@%s>, Please DM satoshi this command instead, the response may contain sensitive information. Thanks", userID)
}

func formatFailureMsg(userID, failureMsg string, err error) string {
	return fmt.Sprintf(
		":disappointed: Sorry <@%s>, I failed to execute that command.\n%s\n Error: %s\n.",
		userID, failureMsg, err,
	)
}

// Placeholder
func normalizeContent(content string) string {
	return content
}
