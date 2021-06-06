package satoshi

import (
	"context"
	"fmt"
	"strings"

	discordproto "swallowtail/s.discord/proto"
	"swallowtail/s.satoshi/coins"
	"swallowtail/s.satoshi/pricebot"

	"github.com/bwmarrin/discordgo"
	"github.com/monzo/slog"
)

const (
	priceCommandID       = "price-comamnd"
	priceCommandUsageMsg = ":wave: <@%s>: `Usage: !price [symbols... | all ]`\nExamples:\n`!price BTC ETH LTC`\n`!price all\n`"
)

var (
	priceBotSvc pricebot.PriceBotService
)

func init() {
	registerSatoshiCommand(priceCommandID, priceCommand)
	priceBotSvc = pricebot.NewService(context.Background())
}

func priceCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Check if we're parsing the correct command
	if !strings.HasPrefix(m.Content, "!price") {
		return
	}
	tokens := strings.Split(m.Content, " ")
	if len(tokens) < 2 {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(priceCommandUsageMsg, m.Author.ID))
		return
	}
	var (
		symbols      []string
		channelID    = m.ChannelID
		withGreeting bool
	)
	// Handle the case for all; set channel to PriceBot channel if true.
	if strings.ToLower(tokens[1]) == "all" {
		symbols = coins.GetCoinSymbols()
		channelID = discordproto.DiscordSatoshiPriceBotChannel
		withGreeting = true
	}
	symbols = tokens[1:]

	pricesMsg := priceBotSvc.GetPricesAsFormattedString(nil, symbols, withGreeting)
	_, err := s.ChannelMessageSend(channelID, pricesMsg)
	if err != nil {
		// Best effort
		slog.Info(nil, "Price command failed to send prices")
	}
}
