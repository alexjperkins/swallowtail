package satoshi

import (
	"context"
	"fmt"
	"strings"
	"time"

	accountproto "swallowtail/s.account/proto"

	"github.com/bwmarrin/discordgo"
	"github.com/monzo/slog"
	"google.golang.org/grpc"
)

const (
	addExchangeCommandID     = "add-exchange-command"
	addExchangeCommandPrefix = "!exchange-add"
)

func init() {
	registerSatoshiCommand(addExchangeCommandID, addExchangeCommand)
}

func addExchangeCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !strings.HasPrefix(m.Content, addExchangeCommandPrefix) {
		return
	}

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(30*time.Second))
	defer cancel()

	tokens := strings.Split(m.Content, " ")
	if len(tokens) != 3 {
		s.ChannelMessageSend(
			m.ChannelID, fmt.Sprintf(
				":wave: Hi, incorrect usage.\n\nExample: `!exchange-add <exchange> <api-key> <secret-key>`\nSupported exchanges `%s, %s`",
				accountproto.ExchangeType_BINANCE.String(),
				accountproto.ExchangeType_FTX.String(),
			),
		)
		return
	}

	slog.Debug(ctx, "Received %s command, args: %v", addExchangeCommandID, tokens)

	exchange, apiKey, secretKey := strings.ToUpper(tokens[1]), tokens[2], tokens[3]
	var exchangeType accountproto.ExchangeType

	switch exchange {
	case accountproto.ExchangeType_BINANCE.String():
		exchangeType = accountproto.ExchangeType_BINANCE
	case accountproto.ExchangeType_FTX.String():
		exchangeType = accountproto.ExchangeType_FTX
	default:
		// Bad Exchange type.
		s.ChannelMessageSend(
			m.ChannelID,
			fmt.Sprintf(":wave: Sorry, I don't support that exchange\n\nPlease choose from `%s, %s`",
				accountproto.ExchangeType_BINANCE.String(),
				accountproto.ExchangeType_FTX.String(),
			),
		)
		return
	}

	conn, err := grpc.DialContext(ctx, "swallowtail-s-account:8000", grpc.WithInsecure())
	if err != nil {
		slog.Error(ctx, "Failed to reach s_account grpc: %v", err)
		return
	}
	defer conn.Close()

	client := accountproto.NewAccountClient(conn)
	_, err = client.AddExchange(ctx, &accountproto.AddExchangeRequest{
		Exchange: &accountproto.Exchange{
			UserId:    m.Author.ID,
			Exchange:  exchangeType,
			ApiKey:    apiKey,
			SecretKey: secretKey,
			IsActive:  true,
		},
	})
	if err != nil {
		slog.Error(ctx, "Failed to add a new exchange to account: %v", err.Error(), map[string]string{
			"user_id":       m.Author.ID,
			"username":      m.Author.Username,
			"exchange":      exchange,
			"exchange_type": exchangeType.String(),
		})
		s.ChannelMessageSend(
			m.ChannelID,
			fmt.Sprintf(":wave: Sorry, I wasn't able to to add an exchange; please ping @ajperkins to investigate."))
		return
	}

	s.ChannelMessageSend(
		m.ChannelID,
		fmt.Sprintf(":wave: Thanks! I've now added the exchange to your account. \n To see all exchanges registered use the command: `!exchange-list`"),
	)
}
