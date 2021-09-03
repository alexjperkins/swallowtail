package commands

import (
	"context"
	"fmt"
	"strings"
	"time"

	"swallowtail/libraries/gerrors"
	accountproto "swallowtail/s.account/proto"

	"google.golang.org/grpc"

	"github.com/bwmarrin/discordgo"
	"github.com/monzo/slog"
)

const (
	exchangeCommandID = "exchange"
	exchangeUsage     = `
	Usage: !exchange <subcommand>

	Subcommands:
	1. register <exchange> <api-key> <secret-key>
	2. list
	`
)

func init() {
	register(exchangeCommandID, &Command{
		ID:                  exchangeCommandID,
		IsPrivate:           true,
		MinimumNumberOfArgs: 1,
		Usage:               exchangeUsage,
		Handler:             exchangeCommand,
		SubCommands: map[string]*Command{
			"register": {
				ID:                  "exchange-register",
				IsPrivate:           true,
				IsFuturesOnly:       true,
				MinimumNumberOfArgs: 3,
				Usage:               `Usage: !exchange register binance <api-key> <secret-key>`,
				Handler:             registerExchangeCommand,
			},
			"list": {
				ID:                  "exchange-list",
				IsPrivate:           true,
				IsFuturesOnly:       true,
				MinimumNumberOfArgs: 0,
				Usage:               `Usage: !exchange list`,
				Handler:             listExchangeCommand,
			},
		},
	})
}

func exchangeCommand(ctx context.Context, tokens []string, s *discordgo.Session, m *discordgo.MessageCreate) error {
	return gerrors.Unimplemented("parent_command_unimplemented.exchange", nil)
}

func registerExchangeCommand(ctx context.Context, tokens []string, s *discordgo.Session, m *discordgo.MessageCreate) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(30*time.Second))
	defer cancel()

	exchange, apiKey, secretKey := strings.ToUpper(tokens[0]), tokens[1], tokens[2]
	var exchangeType accountproto.ExchangeType

	switch exchange {
	case accountproto.ExchangeType_BINANCE.String():
		exchangeType = accountproto.ExchangeType_BINANCE
	case accountproto.ExchangeType_FTX.String():
		exchangeType = accountproto.ExchangeType_FTX
	default:
		// Bad Exchange type.
		if _, err := (s.ChannelMessageSend(
			m.ChannelID,
			fmt.Sprintf(":wave: Sorry, I don't support that exchange\n\nPlease choose from `%s, %s`",
				accountproto.ExchangeType_BINANCE.String(),
				accountproto.ExchangeType_FTX.String(),
			),
		)); err != nil {
			return gerrors.Augment(err, "failed_to_send_to_discord_bad_exchange", map[string]string{
				"command_id": "register-exchange-command",
			})
		}

		return nil
	}

	if _, err := (&accountproto.AddExchangeRequest{
		UserId: m.Author.ID,
		Exchange: &accountproto.Exchange{
			ExchangeType: exchangeType,
			ApiKey:       apiKey,
			SecretKey:    secretKey,
			IsActive:     true,
		},
	}).Send(ctx).Response(); err != nil {
		s.ChannelMessageSend(
			m.ChannelID,
			fmt.Sprintf(":wave: Sorry, I wasn't able to to add an exchange; please ping @ajperkins to investigate."),
		)
		return gerrors.Augment(err, "failed_to_send_to_discord_failure", nil)
	}

	_, err := s.ChannelMessageSend(
		m.ChannelID,
		fmt.Sprintf(":wave: Thanks! I've now added the exchange to your account. \n\n To see all exchanges registered use the command: ```!exchange list```"),
	)

	return err
}

func listExchangeCommand(ctx context.Context, tokens []string, s *discordgo.Session, m *discordgo.MessageCreate) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(30*time.Second))
	defer cancel()

	conn, err := grpc.DialContext(ctx, "swallowtail-s-account:8000", grpc.WithInsecure())
	if err != nil {
		slog.Error(ctx, "Failed to reach s_account grpc: %v", err)
		return err
	}
	defer conn.Close()

	client := accountproto.NewAccountClient(conn)
	rsp, err := client.ListExchanges(ctx, &accountproto.ListExchangesRequest{
		UserId:     m.Author.ID,
		ActiveOnly: true,
	})

	exchanges := rsp.GetExchanges()
	if exchanges == nil {
		_, err := s.ChannelMessageSend(
			m.ChannelID,
			fmt.Sprintf(":wave: Sorry, you don't have any exchanges registered I'm afraid."),
		)
		return err
	}

	exchangesMsg := formatExchangesToMsg(exchanges, m)

	_, err = s.ChannelMessageSend(
		m.ChannelID,
		fmt.Sprintf(":wave: Here's the exchange details registered to your account, all keys are masked\n\n%s", exchangesMsg),
	)

	return err
}
