package commands

import (
	"context"
	"fmt"

	"swallowtail/libraries/gerrors"
	accountproto "swallowtail/s.account/proto"

	"github.com/bwmarrin/discordgo"
	"github.com/monzo/slog"
	"github.com/monzo/terrors"
)

const (
	accountCommandID = "account"
	accountUsage     = `
	Usage: !account <subcommand>

	SubCommands:
	1. register
	`
)

func init() {
	register(accountCommandID, &Command{
		ID:                  accountCommandID,
		Private:             true,
		MinimumNumberOfArgs: 1,
		Usage:               accountUsage,
		Handler:             accountHandler,
		SubCommands: map[string]*Command{
			"register": {
				ID:                  "register-account",
				Private:             true,
				MinimumNumberOfArgs: 2,
				Usage:               `!account register <username> <password>`,
				Handler:             registerAccountHandler,
				FailureMsg:          "Please check you already don't have an account; ping @ajperkins with your message if you need help",
			},
		},
	})
}

func accountHandler(ctx context.Context, tokens []string, s *discordgo.Session, m *discordgo.MessageCreate) error {
	return gerrors.Unimplemented("parent_command_unimplemented.account", nil)
}

func registerAccountHandler(ctx context.Context, tokens []string, s *discordgo.Session, m *discordgo.MessageCreate) error {
	username, password := tokens[0], tokens[1]

	_, err := (&accountproto.CreateAccountRequest{
		UserId:   m.Author.ID,
		Username: username,
		Password: password,
	}).Send(ctx).Response()
	switch {
	case terrors.Is(err, terrors.ErrPreconditionFailed, "account-already-exists"):
		s.ChannelMessageSend(
			m.ChannelID,
			":wave: Hi, I've already got an account registered for you.  You're all good!",
		)
	case err != nil:
		slog.Error(ctx, "Failed to create new account: %v", err, map[string]string{
			"user_id":  m.Author.ID,
			"username": username,
		})
		s.ChannelMessageSend(
			m.ChannelID,
			fmt.Sprintf(":disappointed: Sorry, I failed to create an account with username: `%s`, please ping @ajperkins to investigate. Thanks", username),
		)
		return nil
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":wave: I have registered your account with username: `%s`", username))
	return nil
}
