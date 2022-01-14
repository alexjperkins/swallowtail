package commands

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/monzo/slog"
	"google.golang.org/grpc"

	"swallowtail/libraries/gerrors"
	accountproto "swallowtail/s.account/proto"
	tradeengineproto "swallowtail/s.trade-engine/proto"
)

const (
	exchangeCommandID = "exchange"
	exchangeUsage     = `!exchange <subcommand>`
)

func init() {
	register(exchangeCommandID, &Command{
		ID:                  exchangeCommandID,
		IsPrivate:           true,
		IsFuturesOnly:       true,
		MinimumNumberOfArgs: 1,
		Usage:               exchangeUsage,
		Handler:             exchangeCommand,
		Description:         "Manages all things related to exchanges such as api keys & more.",
		Guide:               "https://scalloped-single-1bd.notion.site/How-to-register-an-exchange-d3d73af635f041a89a3e57d3d33a32b0",
		SubCommands: map[string]*Command{
			"register": {
				ID:                  "exchange-register",
				IsPrivate:           true,
				IsFuturesOnly:       true,
				MinimumNumberOfArgs: 3,
				Usage:               `!exchange register binance <api-key> <secret-key>`,
				Description:         "Registers a set of API keys (Binance only for now).",
				Handler:             registerExchangeCommand,
			},
			"list": {
				ID:                  "exchange-list",
				IsPrivate:           true,
				IsFuturesOnly:       true,
				MinimumNumberOfArgs: 0,
				Usage:               `!exchange list`,
				Description:         "Lists all registered API keys.",
				Handler:             listExchangeCommand,
			},
			"set-primary": {
				ID:                  "exchange-set-primary",
				IsPrivate:           true,
				IsFuturesOnly:       true,
				MinimumNumberOfArgs: 1,
				Usage:               `set-primary <exchange>`,
				Description:         "Sets the primary exchange to use on your account",
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

	venue, apiKey, secretKey := strings.ToUpper(tokens[0]), tokens[1], tokens[2]
	var venueProto tradeengineproto.VENUE

	switch strings.ToUpper(venue) {
	case tradeengineproto.VENUE_BINANCE.String():
		venueProto = tradeengineproto.VENUE_BINANCE
	case tradeengineproto.VENUE_BITFINEX.String():
		venueProto = tradeengineproto.VENUE_BITFINEX
	case tradeengineproto.VENUE_DERIBIT.String():
		venueProto = tradeengineproto.VENUE_DERIBIT
	case tradeengineproto.VENUE_FTX.String():
		venueProto = tradeengineproto.VENUE_FTX
	default:
		// Bad Exchange type.
		if _, err := (s.ChannelMessageSend(
			m.ChannelID,
			fmt.Sprintf(":wave: Sorry, I don't support that venues\n\nPlease post in #crypto-support to check available venues`%s, %s`"),
		)); err != nil {
			return gerrors.Augment(err, "failed_to_send_to_discord_bad_exchange", map[string]string{
				"command_id": "register-exchange-command",
			})
		}

		return nil
	}

	rsp, err := (&accountproto.AddVenueAccountRequest{
		UserId: m.Author.ID,
		VenueAccount: &accountproto.VenueAccount{
			Venue:     venueProto,
			ApiKey:    apiKey,
			SecretKey: secretKey,
			IsActive:  true,
		},
	}).Send(ctx).Response()
	if err != nil {
		slog.Error(ctx, "Failed to add exchange, error: %v", err)
		_, derr := s.ChannelMessageSend(
			m.ChannelID,
			fmt.Sprintf(":wave: Sorry, I wasn't able to to add an exchange; please ping @ajperkins to investigate."),
		)
		if derr != nil {
			return gerrors.Augment(derr, "failed_to_send_to_discord_failure", nil)
		}

		return nil
	}

	if !rsp.Verified {
		// Convert reasons into human friendly format.
		var reasons strings.Builder
		for _, r := range strings.Split(rsp.Reason, ",") {
			reasons.WriteString(fmt.Sprintf("- %s\n", r))
		}

		_, derr := s.ChannelMessageSend(
			m.ChannelID,
			fmt.Sprintf(
				":wave: Sorry, I wasn't able to to verify your credentials. This is likely due to the following permissisions issues:```%s```",
				reasons.String(),
			),
		)
		if derr != nil {
			return gerrors.Augment(derr, "failed_to_send_to_discord_failure", nil)
		}

		return nil
	}

	_, err = s.ChannelMessageSend(
		m.ChannelID,
		fmt.Sprintf(":wave: Thanks! I've now added the exchange to your account. \n\n To see all exchanges registered use the command: ```!exchange list```"),
	)

	return err
}

func listExchangeCommand(ctx context.Context, tokens []string, s *discordgo.Session, m *discordgo.MessageCreate) error {
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(30*time.Second))
	defer cancel()

	conn, err := grpc.DialContext(ctx, "swallowtail-s-account:8000", grpc.WithInsecure())
	if err != nil {
		slog.Error(ctx, "Failed to reach s_account grpc: %v", err)
		return err
	}
	defer conn.Close()

	client := accountproto.NewAccountClient(conn)
	rsp, err := client.ListVenueAccounts(ctx, &accountproto.ListVenueAccountsRequest{
		UserId:     m.Author.ID,
		ActiveOnly: true,
	})

	venueAccounts := rsp.GetVenueAccounts()
	if venueAccounts == nil {
		_, err := s.ChannelMessageSend(
			m.ChannelID,
			fmt.Sprintf(":wave: Sorry, you don't have any exchanges registered I'm afraid."),
		)
		return err
	}

	exchangesMsg := formatVenueAccountsToMsg(venueAccounts, m)

	_, err = s.ChannelMessageSend(
		m.ChannelID,
		fmt.Sprintf(":wave: Here's the exchange details registered to your account, all keys are masked\n\n%s", exchangesMsg),
	)

	return err
}

func setPrimaryExchangeCommand(ctx context.Context, tokens []string, s *discordgo.Session, m *discordgo.MessageCreate) error {
	venueToken := tokens[0]

	var venue tradeengineproto.VENUE
	switch strings.ToUpper(venueToken) {
	case strings.ToUpper(tradeengineproto.VENUE_BINANCE.String()):
		venue = tradeengineproto.VENUE_BINANCE
	case strings.ToUpper(tradeengineproto.VENUE_FTX.String()):
		venue = tradeengineproto.VENUE_FTX
	case strings.ToUpper(tradeengineproto.VENUE_DERIBIT.String()):
		return gerrors.Unimplemented("venue_unimplemented_for_primary_account.deribit", nil)
	case strings.ToUpper(tradeengineproto.VENUE_BITFINEX.String()):
		return gerrors.Unimplemented("venue_unimplemented_for_primary_account.deribit", nil)
	default:
		return gerrors.Unimplemented("venue_unimplemented", map[string]string{
			"venue": venueToken,
		})
	}

	if _, err := (&accountproto.UpdateAccountRequest{
		PrimaryVenue: venue,
	}).SendWithTimeout(ctx, 10*time.Second).Response(); err != nil {
		s.ChannelMessageSend(
			m.ChannelID,
			fmt.Sprintf(
				":wave I was unable to set your primary exchange on your account to: %s, Error: %v", venue, err,
			),
		)

		return gerrors.Augment(err, "failed_to_set_primary_venue", map[string]string{
			"venue": venueToken,
		})
	}

	s.ChannelMessageSend(
		m.ChannelID,
		fmt.Sprintf(
			":wave I have set your primary exchange on your account to be: %s", venue,
		),
	)

	return nil
}
