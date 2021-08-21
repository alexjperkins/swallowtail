package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Command ...
type Command struct {
	ID      string
	Prefix  string
	Private bool
	// Non-inclusive of the prefix.
	MinimumNumberOfArgs int
	Usage               string
	Handler             CommandHandler
	SubCommands         map[string]*SubCommand
}

type SubCommand struct {
	ID                  string
	MinimumNumberOfArgs int
	Usage               string
	Handler             CommandHandler
	SubCommands         map[string]*SubCommand
}

// CommandHandler ...
type CommandHandler func(ctx context.Context, tokens []string, s *discordgo.Session, m *discordgo.MessageCreate) error

// Exec ...
func (c *Command) Exec(s *discordgo.Session, m *discordgo.MessageCreate) error {
	// Check Prefix
	if !strings.HasPrefix(m.Content, c.Prefix) {
		return nil
	}

	// Check Usage
	tokens := strings.Split(m.Content, " ")
	if len(tokens) > 1 {
		if strings.ToLower(tokens[1]) == "help" {
			_, err := s.ChannelMessageSend(m.ChannelID, formatUsageMsg(m.Author.ID, c.Usage))
			return err
		}
	}

	tokens = tokens[1:]
	// Check we have at least the correct number of arguments to execute the command.
	if len(tokens) < c.MinimumNumberOfArgs {
		_, err := s.ChannelMessageSend(m.ChannelID, formatUsageMsg(m.Author.ID, c.Usage))
		return err
	}

	ctx := context.Background()

	// Check if command should be in DMs
	if c.Private && m.GuildID != "" {
		_, err := s.ChannelMessageSend(m.ChannelID, formatNonPublicMsg(m.Author.ID))
		return err
	}

	subCommand, ok := c.SubCommands[tokens[0]]
	if !ok {
		// We don't have a subcommand that matches the "second" token; we can
		// try to run the original command instead.
		if err := c.Handler(ctx, tokens, s, m); err != nil {
			_, err := s.ChannelMessageSend(m.ChannelID, formatFailureMsg(m.Author.ID, err.Error()))
			return err
		}
	}

	tokens = tokens[1:]

	// Check if the next arg is for help.
	if len(tokens) > 1 {
		if strings.ToLower(tokens[1]) == "help" {
			_, err := s.ChannelMessageSend(m.ChannelID, formatUsageMsg(m.Author.ID, c.Usage))
			return err
		}
	}

	// Check we have at least the correct number of arguments; if we do we can execute the subcommand.
	if len(tokens) < subCommand.MinimumNumberOfArgs {
		_, err := s.ChannelMessageSend(m.ChannelID, formatUsageMsg(m.Author.ID, subCommand.Usage))
		return err
	}

	if err := subCommand.Exec(ctx, tokens, s, m); err != nil {
		_, err := s.ChannelMessageSend(m.ChannelID, formatFailureMsg(m.Author.ID, err.Error()))
		return err
	}

	return nil
}

func (c *SubCommand) Exec(ctx context.Context, tokens []string, s *discordgo.Session, m *discordgo.MessageCreate) error {
	subCommand, ok := c.SubCommands[tokens[0]]
	if !ok {
		// We don't have a subcommand that matches the "second" token; we can
		// try to run the original command instead.
		if err := c.Handler(ctx, tokens, s, m); err != nil {
			_, err := s.ChannelMessageSend(m.ChannelID, formatFailureMsg(m.Author.ID, err.Error()))
			return err
		}
	}

	tokens = tokens[1:]
	// Check if the next arg is for help.
	if len(tokens) > 1 {
		if strings.ToLower(tokens[0]) == "help" {
			_, err := s.ChannelMessageSend(m.ChannelID, formatUsageMsg(m.Author.ID, c.Usage))
			return err
		}
	}

	// Check we have at least the correct number of arguments; if we do we can execute the subcommand.
	if len(tokens) < subCommand.MinimumNumberOfArgs {
		_, err := s.ChannelMessageSend(m.ChannelID, formatUsageMsg(m.Author.ID, subCommand.Usage))
		return err
	}

	if err := subCommand.Exec(ctx, tokens, s, m); err != nil {
		_, err := s.ChannelMessageSend(m.ChannelID, formatFailureMsg(m.Author.ID, err.Error()))
		return err
	}

	return nil
}

func formatUsageMsg(userID, usage string) string {
	return fmt.Sprintf(":wave: <@%s> Usage: %s", userID, usage)
}

func formatNonPublicMsg(userID string) string {
	return fmt.Sprintf(":wave: <@%s>, Please DM satoshi this command instead. Thanks", userID)
}

func formatFailureMsg(userID, failureMsg string) string {
	return fmt.Sprintf(
		":disappointed: Sorry <@%s>, I failed to execute that command.\n Error: %s\nPlease check the usage is correct.",
		userID, failureMsg,
	)
}
