package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	helpCommandID    = "help"
	helpCommandUsage = `
	Usage: !help
	`
)

func init() {
	register(helpCommandID, &Command{
		ID:      helpCommandID,
		Usage:   helpCommandUsage,
		Handler: helpCommand,
	})
}

func helpCommand(ctx context.Context, tokens []string, s *discordgo.Session, m *discordgo.MessageCreate) error {
	commands := List()

	var sb strings.Builder
	sb.WriteString("Satoshi Commands: \n")
	for _, command := range commands {
		sb.WriteString(fmt.Sprintf("\n%s\n\n%s", strings.ToTitle(command.ID), command.Usage))
	}

	s.ChannelMessageSend(
		m.ChannelID,
		sb.String(),
	)

	return nil
}
