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

	NOTE: "<...>" means you need to replace with the argument value, don't include the brackets.

	To see what args are required for subcommands, call help as an argument.

	Some commands have guides attached; run "!<command> help" to see the guide if it's there.

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
	for _, command := range commands[:len(commands)/2] {
		sb.WriteString(fmt.Sprintf("\n%s%s", strings.ToTitle(command.ID), command.Usage))
	}

	// TODO: this breaks if we're over 2000 chars
	// This temp fix is super awkward; we want to improve it at some point.
	_, err := s.ChannelMessageSend(
		m.ChannelID,
		util.WrapAsCodeBlock(sb.String()),
	)
	if err != nil {
		return err
	}

	sb.Reset()
	for _, command := range commands[len(commands)/2:] {
		sb.WriteString(fmt.Sprintf("\n%s%s", strings.ToTitle(command.ID), command.Usage))
	}

	_, err = s.ChannelMessageSend(
		m.ChannelID,
		util.WrapAsCodeBlock(sb.String()),
	)
	if err != nil {
		return err
	}

	return err
}
