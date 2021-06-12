package satoshi

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/monzo/slog"
)

const (
	riskCommandID     = "risk-command"
	riskCommandPrefix = "!risk"
)

func init() {
	registerSatoshiCommand(riskCommandID, riskCalculator)
}

func riskCalculator(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !strings.HasPrefix(m.Content, riskCommandPrefix) {
		return
	}

	tokens := strings.Split(m.Content, " ")
	if len(tokens) != 5 {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Hi @%s, `!risk usage: <entry> <stop loss> <account size> <percentage eg 0.05>`", m.Author.Username))
		return
	}
	slog.Info(context.TODO(), "Received %s command, args: %v", riskCommandPrefix, tokens)

	entry, err := strconv.ParseFloat(tokens[1], 64)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Hi @%s, couldn't parse entry: %v into a float, please check.", m.Author.Username, tokens[1]))
		return
	}
	stopLoss, err := strconv.ParseFloat(tokens[2], 64)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Hi @%s, couldn't parse stop loss: %v into a float, please check.", m.Author.Username, tokens[2]))
		return
	}
	accountSize, err := strconv.ParseFloat(tokens[3], 64)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Hi @%s, couldn't parse accountSize: %v into a float, please check.", m.Author.Username, tokens[3]))
		return
	}
	percentage, err := strconv.ParseFloat(tokens[4], 64)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Hi @%s, couldn't parse percentage: %v into a float, please check.", m.Author.Username, tokens[4]))
		return
	}

	contracts := calculateRisk(entry, stopLoss, accountSize, percentage)

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Hi @%s, you need to buy **%.2f** contracts for %v%% risk.", m.Author.Username, contracts, percentage*100))
	return
}
