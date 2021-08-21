package commands

import (
	"context"
	"fmt"

	accountproto "swallowtail/s.account/proto"

	"github.com/bwmarrin/discordgo"
	"github.com/monzo/slog"
	"github.com/monzo/terrors"
	"google.golang.org/grpc"
)

const (
	accountCommandID     = "account-command"
	accountCommandPrefix = "!account"
	accountUsage         = `
	Usage: !account <subcommand>

	SubCommands:
	1. register
	`
)

func init() {
	register(accountCommandID, &Command{
		ID:                  accountCommandID,
		Prefix:              accountCommandPrefix,
		Private:             true,
		MinimumNumberOfArgs: 1,
		Usage:               accountUsage,
		Handler:             accountHandler,
		SubCommands: map[string]*SubCommand{
			"register": {
				ID:                  "register-account-command",
				MinimumNumberOfArgs: 2,
				Usage:               `!account register <username> <password>`,
				Handler:             registerAccountHandler,
			},
		},
	})
}

func accountHandler(ctx context.Context, tokens []string, s *discordgo.Session, m *discordgo.MessageCreate) error {
	return nil
}

func registerAccountHandler(ctx context.Context, tokens []string, s *discordgo.Session, m *discordgo.MessageCreate) error {
	username, password := tokens[0], tokens[1]

	conn, err := grpc.DialContext(ctx, "swallowtail-s-account:8000", grpc.WithInsecure())
	if err != nil {
		slog.Error(ctx, "Failed to reach s_account grpc: %v", err)
		return nil
	}
	defer conn.Close()

	client := accountproto.NewAccountClient(conn)
	_, err = client.CreateAccount(ctx, &accountproto.CreateAccountRequest{
		UserId:   m.Author.ID,
		Username: username,
		Password: password,
	})
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

	slog.Info(ctx, "Created new account: %s: %s", m.Author.Username, m.Author.ID)

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":wave: I have registered your account with username: `%s`", username))
	return nil
}
