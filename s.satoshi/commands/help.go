package commands

import (
	"context"
	"strings"
	"swallowtail/libraries/gerrors"
	"swallowtail/libraries/util"

	"github.com/bwmarrin/discordgo"
)

const (
	helpCommandID    = "help"
	helpCommandUsage = `!help`

	helpMessage = `
	Satoshi works by parsing commands prefixed with a '!' idenitifier.
	To call a command simply message satoshi the command.

	Some commands also have subcommands. Some commands can only be ran by certain memmbers.

	Some commands & subcommands also require "args", these are values given to satoshi
	as part of the command or subcommand.

	NOTE: "<...>" means you need to replace with the argument value, don't include the brackets.

	To see what args are required for subcommands, call help as an argument.

	Please ping @ajperkins if you have any questions.`
)

func init() {
	register(helpCommandID, &Command{
		ID:          helpCommandID,
		Usage:       helpCommandUsage,
		Handler:     helpCommand,
		Description: "Prints all available commands and subcommands",
	})
}

func helpCommand(ctx context.Context, tokens []string, s *discordgo.Session, m *discordgo.MessageCreate) error {
	commands := List()

	membersRoles, err := getMembersRolesFromGuild(s, m.Author.ID)
	if err != nil {
		return gerrors.Augment(err, "help_command_failed.failed_to_get_guild_member_roles", nil)
	}

	futuresMember := isFuturesMember(membersRoles)
	admin := isAdmin(membersRoles)

	var sb strings.Builder

	sb.WriteString("SATOSHI COMMANDS")
	sb.WriteString(helpMessage)

	for _, command := range commands[:len(commands)/2] {
		sb.WriteString(formatHelpMsg(command, futuresMember, admin))
	}

	// TODO: this breaks if we're over 2000 chars
	// This temp fix is super awkward; we want to improve it at some point.
	_, err = s.ChannelMessageSend(
		m.ChannelID,
		util.WrapAsCodeBlock(sb.String()),
	)
	if err != nil {
		return err
	}

	sb.Reset()
	for _, command := range commands[len(commands)/2:] {
		sb.WriteString(formatHelpMsg(command, futuresMember, admin))
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
