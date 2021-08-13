package satoshi

import (
	"context"
	"fmt"
	"strings"

	accountproto "swallowtail/s.account/proto"

	"github.com/bwmarrin/discordgo"
	"github.com/monzo/slog"
	"github.com/monzo/terrors"
)

const (
	exchangeCommandID     = "exchange-command"
	exchangeCommandPrefix = "!exchange"
	exchangeCommandUsage  = `
	Usage: !exchange <operation> <args>
	Example: !exchange binance this-is-an-api-key this-is-a-secret-key

	Operations:
	1. register <exchange> <api-key> <secret-key>

	Exchanges:
	1. binance
	`
)

func init() {
	registerSatoshiCommand(exchangeCommandID, exchangeCommand)
}

func exchangeCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !strings.HasPrefix(m.Content, exchangeCommandPrefix) {
		return
	}

	ctx := context.Background()

	tokens := strings.Split(m.Content, " ")
	if len(tokens) < 2 {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":wave: <@%s> \n`%s`", m.Author.ID, exchangeCommandUsage))
	}

	switch strings.ToLower(tokens[1]) {
	case "register":
		if len(tokens[2:]) < 3 {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":wave: <@%s> \n`%s`", m.Author.ID, exchangeCommandUsage))
			return
		}
		e, ak, sk := tokens[2], tokens[3], tokens[4]
		if err := registerExchange(ctx, m.Author.ID, e, ak, sk); err != nil {
			slog.Error(ctx, "Failed to register exchange: %v", err)

			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":wave: <@%s> \nSorry, I failed to regsiter that exchange, please check the keys are correct\nError: %s", m.Author.ID, err))
			return
		}

		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":wave: <@%s>, exchange registered!", m.Author.ID))
		return
	}
}

func registerExchange(ctx context.Context, userID, exchange, apiKey, secretKey string) error {
	if _, err := (&accountproto.AddExchangeRequest{
		Exchange: &accountproto.Exchange{
			UserId:    userID,
			ApiKey:    apiKey,
			SecretKey: secretKey,
			IsActive:  true,
		},
	}).Send(ctx).Response(); err != nil {
		return terrors.Augment(err, "Failed to register exchange", nil)
	}
	return nil
}
