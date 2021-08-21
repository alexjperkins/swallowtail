package commands

import (
	"context"
	"fmt"
	"strings"
	"swallowtail/libraries/util"

	"github.com/bwmarrin/discordgo"
)

const (
	helpCommandID    = "help"
	helpCommandUsage = `
	Usage: !help
	Description: prints all available commands and subcommands.
	`

	helpMessage = `
	Satoshi works by parsing commands prefixed with a '!' idenitifier.
	To call a command simply message satoshi the command.

	Some commands also have subcommands. Subcommands don't require the identifier,
	they can just follow the command.

	Some commands & subcommands also require "args", these are values given to satoshi
	as part of the command or subcommand.

	To see what args are required for subcommands, call help as an argument.
	See notion for more detail.

	Please ping @ajperkins if you have any questions.
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
	sb.WriteString("SATOSHI COMMANDS")
	sb.WriteString(helpMessage)
	for _, command := range commands {
		sb.WriteString(fmt.Sprintf("\n%s%s", strings.ToTitle(command.ID), command.Usage))
	}

	s.ChannelMessageSend(
		m.ChannelID,
		util.WrapAsCodeBlock(sb.String()),
	)

	return nil
}
