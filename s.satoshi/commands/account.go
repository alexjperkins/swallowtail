package commands

import (
	"context"
	"fmt"

	"swallowtail/libraries/gerrors"
	"swallowtail/libraries/util"
	accountproto "swallowtail/s.account/proto"

	"github.com/bwmarrin/discordgo"
	"github.com/monzo/slog"
	"github.com/monzo/terrors"
)

const (
	accountCommandID = "account"
	accountUsage     = `!account <subcommand>`
)

func init() {
	register(accountCommandID, &Command{
		ID:                  accountCommandID,
		IsPrivate:           true,
		MinimumNumberOfArgs: 1,
		Usage:               accountUsage,
		Description:         "Command for managing satoshi account",
		Handler:             accountHandler,
		SubCommands: map[string]*Command{
			"register": {
				ID:                  "account-register",
				IsPrivate:           true,
				MinimumNumberOfArgs: 2,
				Usage:               `!account register <email> <password>`,
				Description:         "Manages everything related to your account.",
				Handler:             registerAccountHandler,
				FailureMsg:          "Please check you already don't have an account; ping @ajperkins with your message if you need help",
			},
			"read": {
				ID:                  "account-read",
				IsPrivate:           true,
				MinimumNumberOfArgs: 0,
				Usage:               `!account read`,
				Description:         "Returns everything satoshi stores as your account. You can see if you have an account with this.",
				Handler:             readAccountHandler,
			},
		},
	})
}

func accountHandler(ctx context.Context, tokens []string, s *discordgo.Session, m *discordgo.MessageCreate) error {
	return gerrors.Unimplemented("parent_command_unimplemented.account", nil)
}

func registerAccountHandler(ctx context.Context, tokens []string, s *discordgo.Session, m *discordgo.MessageCreate) error {
	email, password := tokens[0], tokens[1]

	_, err := (&accountproto.CreateAccountRequest{
		UserId:   m.Author.ID,
		Username: m.Author.Username,
		Email:    email,
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
			"user_id": m.Author.ID,
			"email":   email,
		})
		s.ChannelMessageSend(
			m.ChannelID,
			fmt.Sprintf(":disappointed: Sorry, I failed to create an account with email: `%s`, please ping @ajperkins to investigate. Thanks", email),
		)
		return nil
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":wave: I have registered your account with email: `%s`", email))
	return nil
}

func readAccountHandler(ctx context.Context, tokens []string, s *discordgo.Session, m *discordgo.MessageCreate) error {
	rsp, err := (&accountproto.ReadAccountRequest{
		UserId: m.Author.ID,
	}).Send(ctx).Response()
	switch {
	case gerrors.Is(err, gerrors.ErrNotFound):
		s.ChannelMessageSend(
			m.ChannelID,
			":disappointed: Looks like you haven't registered an account with us just yet! Use `!account register help` for help.",
		)
		return nil
	case err != nil:
		return gerrors.Augment(err, "failed_to_read_account", nil)
	case rsp.GetAccount() == nil:
		return gerrors.NotFound("failed_to_read_account.nil_account", map[string]string{
			"user_id": m.Author.ID,
		})
	}

	account := rsp.GetAccount()

	tpl := `
Username:          %s
Email:             %s
Created:           %s
Last Updated:      %v
Is Futures Member: %v
Primary Exchange:  %s
	`
	formattedMsg := fmt.Sprintf(tpl, account.Username, account.Email, account.Created.AsTime(), account.LastUpdated.AsTime(), account.IsFuturesMember, account.PrimaryExchange)

	// Best Effort.
	s.ChannelMessageSend(
		m.ChannelID,
		fmt.Sprintf(":wave: <@%s> Here's your account: %s", m.Author.ID, util.WrapAsCodeBlock(formattedMsg)),
	)

	return nil
}
