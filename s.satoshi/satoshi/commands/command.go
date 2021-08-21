package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Command ...
type Command struct {
	ID                  string
	Prefix              string
	DMOnly              bool
	MinimumNumberOfArgs int
	Usage               string
	exec                func(tokens []string, s *discordgo.Session, m *discordgo.MessageCreate) error
	SubCommands         map[string]*Command
}

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

	if len(tokens) < c.MinimumNumberOfArgs {
		_, err := s.ChannelMessageSend(m.ChannelID, formatUsageMsg(m.Author.ID, c.Usage))
		return err
	}

	// Check if command should be in DMs
	if c.DMOnly && m.GuildID != "" {
		_, err := s.ChannelMessageSend(m.ChannelID, formatNonPublicMsg(m.Author.ID))
		return err
	}

	// Execute command
	if err := c.exec(tokens, s, m); err != nil {
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
