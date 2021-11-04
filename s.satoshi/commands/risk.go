package commands

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"swallowtail/libraries/gerrors"

	"github.com/bwmarrin/discordgo"
)

const (
	riskCommandID    = "risk"
	riskCommandUsage = `!risk <entry> <stop loss> <account size> <percentage of account>`
)

func init() {
	register(riskCommandID, &Command{
		ID:                  riskCommandID,
		Usage:               riskCommandUsage,
		Description:         `A risk calculator; determine how many contracts to buy / sell`,
		MinimumNumberOfArgs: 4,
		Handler:             riskCalculator,
	})
}

func riskCalculator(ctx context.Context, tokens []string, s *discordgo.Session, m *discordgo.MessageCreate) error {
	entry, err := strconv.ParseFloat(tokens[0], 64)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Hi @%s, couldn't parse entry: %v into a float, please check.", m.Author.Username, tokens[1]))
		return gerrors.Augment(err, "bad_param.failed_to_parse.entry", nil)
	}
	stopLoss, err := strconv.ParseFloat(tokens[1], 64)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Hi @%s, couldn't parse stop loss: %v into a float, please check.", m.Author.Username, tokens[2]))
		return gerrors.Augment(err, "bad_param.failed_to_parse.stop_loss", nil)
	}
	accountSize, err := strconv.ParseFloat(tokens[2], 64)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Hi @%s, couldn't parse accountSize: %v into a float, please check.", m.Author.Username, tokens[3]))
		return gerrors.Augment(err, "bad_param.failed_to_parse.account_size", nil)
	}
	percentage, err := strconv.ParseFloat(strings.ReplaceAll(tokens[3], "%", ""), 64)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Hi @%s, couldn't parse percentage: %v into a float, please check.", m.Author.Username, tokens[4]))
		return gerrors.Augment(err, "bad_param.failed_to_parse.percentage", nil)
	}

	contracts := calculateRisk(entry, stopLoss, accountSize, percentage)

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Hi @%s, you need to buy **%.2f** contracts for %v%% risk.", m.Author.Username, contracts, percentage*100))
	return nil
}
