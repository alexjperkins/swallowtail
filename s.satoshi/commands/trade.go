package commands

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"swallowtail/libraries/gerrors"
	tradeengineproto "swallowtail/s.trade-engine/proto"

	"github.com/bwmarrin/discordgo"
	"github.com/monzo/slog"
)

const (
	tradeCommandID = "trade"
	tradeUsage     = `!trade <subcommand>`
)

func init() {
	register(tradeCommandID, &Command{
		ID:                  tradeCommandID,
		IsPrivate:           false,
		MinimumNumberOfArgs: 1,
		Usage:               tradeUsage,
		Description:         "Command for managing satoshi trades",
		Handler:             tradeHandler,
		Guide:               "https://scalloped-single-1bd.notion.site/Automated-trades-guide-69188bbfbb2f4b6f97b22d5cc2a5ee9e",
		SubCommands: map[string]*Command{
			"execute": {
				ID:                  "trade-execute",
				IsPrivate:           false,
				IsFuturesOnly:       true,
				MinimumNumberOfArgs: 2,
				Usage:               `!trade execute <trade_id> <risk (%)>`,
				Handler:             executeTradeHandler,
				FailureMsg:          "Please check the guide you have do the command correctly. Run `!trade help` to see it.",
			},
		},
	})
}

func tradeHandler(ctx context.Context, tokens []string, s *discordgo.Session, m *discordgo.MessageCreate) error {
	return gerrors.Unimplemented("parent_command_unimplemented.trade", nil)
}

func executeTradeHandler(ctx context.Context, tokens []string, s *discordgo.Session, m *discordgo.MessageCreate) error {
	tradeID, riskStr := tokens[0], tokens[1]

	errParams := map[string]string{
		"risk":     riskStr,
		"trade_id": tradeID,
		"user_id":  m.Author.ID,
	}

	// Convert risk to a float
	strings.ReplaceAll(riskStr, "%", "")
	risk, err := strconv.ParseFloat(riskStr, 32)
	if err != nil {
		return gerrors.Augment(err, "failed_to_execute_trade.invalid_risk", errParams)
	}

	if _, err := (&tradeengineproto.AddParticipantToTradeRequest{
		ActorId: tradeengineproto.TradeEngineActorSatoshiSystem,
		UserId:  m.Author.ID,
		TradeId: tradeID,
		IsBot:   false,
		Risk:    float32(risk),
	}).Send(ctx).Response(); err != nil {
		_, err = s.ChannelMessageSend(
			m.ChannelID,
			fmt.Sprintf(":wave:<@%s>, very sorry but it seems as though the trade failed! Error: %v", m.Author.ID, err),
		)
		if err != nil {
			slog.Error(ctx, err.Error())
		}

		return nil
	}

	_, err = s.ChannelMessageSend(
		m.ChannelID,
		fmt.Sprintf(":wave: <@%s> Trade has been successfully executed with %.2f risk :rocket:. Please check manually that everything is in order! :coin:", m.Author.ID, risk),
	)
	if err != nil {
		slog.Error(ctx, err.Error())
	}

	return nil
}
