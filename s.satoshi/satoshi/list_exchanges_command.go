package satoshi

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/monzo/slog"
	"google.golang.org/grpc"

	"swallowtail/libraries/util"
	accountproto "swallowtail/s.account/proto"
)

const (
	listExchangeCommandID     = "list-exchange-command"
	listExchangeCommandPrefix = "!exchange-list"
)

func init() {
	registerSatoshiCommand(listExchangeCommandID, listExchangeCommand)
}

func listExchangeCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !strings.HasPrefix(m.Content, listExchangeCommandPrefix) {
		return
	}

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(30*time.Second))
	defer cancel()

	slog.Debug(ctx, "Received %s command, args: %v", listExchangeCommandID, m.Author.Username)

	conn, err := grpc.DialContext(ctx, "swallowtail-s-account:8000", grpc.WithInsecure())
	if err != nil {
		slog.Error(ctx, "Failed to reach s_account grpc: %v", err)
		return
	}
	defer conn.Close()

	client := accountproto.NewAccountClient(conn)
	rsp, err := client.ListExchanges(ctx, &accountproto.ListExchangesRequest{
		UserId:     m.Author.ID,
		ActiveOnly: true,
	})

	exchanges := rsp.Exchanges
	exchangesMsg := formatExchangesToMsg(exchanges, m)

	s.ChannelMessageSend(
		m.ChannelID,
		fmt.Sprintf(":wave: Here's the exchange details registered to your account, all keys are masked\n\n%s", exchangesMsg),
	)
}

func formatExchangesToMsg(exchanges []*accountproto.Exchange, m *discordgo.MessageCreate) string {
	var lines = []string{}
	lines = append(lines, "`Exchange: ID Username MaskedAPIKey MaskedSecretKey`")
	for i, exchange := range exchanges {
		// We're masking here to be on the safe side; we should expect them to already be masked.
		// TODO maybe we should ping someone here or something.
		maskedAPIKey, maskedSecretKey := util.MaskKey(exchange.ApiKey, 4), util.MaskKey(exchange.SecretKey, 4)

		line := fmt.Sprintf(
			"`%v) %s: %s %s %s %s`",
			i, exchange.Exchange, exchange.ExchangeId, m.Author.Username, maskedAPIKey, maskedSecretKey,
		)
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}
