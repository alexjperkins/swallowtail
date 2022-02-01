package commands

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"

	"swallowtail/libraries/gerrors"
	"swallowtail/libraries/util"
	accountproto "swallowtail/s.account/proto"
)

const (
	accountCommandID = "account"
	accountUsage     = `!account <subcommand>`
)

const (
	binanceReferralLink = "https://www.binance.com/en/futures/ref/swallowtail"
	ftxReferralLink     = "https://ftx.com/#a=9169159"
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
				MinimumNumberOfArgs: 0,
				Usage:               `!account register`,
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
	_, err := (&accountproto.CreateAccountRequest{
		UserId:   m.Author.ID,
		Username: m.Author.Username,
		Email:    m.Author.Email,
	}).Send(ctx).Response()
	switch {
	case gerrors.Is(err, gerrors.ErrAlreadyExists, "account-already-exists"):
		s.ChannelMessageSend(
			m.ChannelID,
			":wave: Hi, I've already got an account registered for you.  You're all good!",
		)
		return nil
	case err != nil:
		return gerrors.Augment(err, "failed_register_account", nil)
	}

	s.ChannelMessageSend(
		m.ChannelID,
		fmt.Sprintf(
			":wave: I have registered your account with email: `%s`.\n\nRef Links:\n`Binance`: %s\n`FTX`: %s",
			m.Author.Email,
			binanceReferralLink,
			ftxReferralLink,
		),
	)
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
Username:               %s
Email:                  %s
Created:                %s
Last Updated:           %v
Is Futures Member:      %v
Primary Venue:          %v
	`
	formattedMsg := fmt.Sprintf(tpl, account.Username, account.Email, account.Created.AsTime(), account.LastUpdated.AsTime(), account.IsFuturesMember, account.PrimaryVenue)

	// Best Effort.
	s.ChannelMessageSend(
		m.ChannelID,
		fmt.Sprintf(":wave: <@%s> Here's your account: %s", m.Author.ID, util.WrapAsCodeBlock(formattedMsg)),
	)

	return nil
}
