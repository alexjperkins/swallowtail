package commands

import (
	"context"
	"time"

	discordproto "swallowtail/s.discord/proto"
	"swallowtail/s.satoshi/coins"
	"swallowtail/s.satoshi/pricebot"

	"github.com/bwmarrin/discordgo"
)

const (
	priceCommandID    = "price"
	priceCommandUsage = `
	Usage: !price <[symbols... | all ]>
	Example: !price BTC ETH LINK
	Description: price command fetches the latest price from coingecko for the symbols provided.

	SubCommands:
	1. all: returns the price of all coins that are registered as default to satoshi.
	`
)

var (
	priceBotSvc pricebot.PriceBotService
)

func init() {
	register(priceCommandID, &Command{
		ID:      priceCommandID,
		Usage:   priceCommandUsage,
		Handler: priceCommand,
		SubCommands: map[string]*Command{
			"all": {
				ID:                  "price-all",
				MinimumNumberOfArgs: 0,
				Usage:               `!price all`,
				Handler:             allPriceCommand,
			},
		},
	})
	priceBotSvc = pricebot.NewService(context.Background())
}

func priceCommand(ctx context.Context, tokens []string, s *discordgo.Session, m *discordgo.MessageCreate) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	symbols := tokens
	pricesMsg := priceBotSvc.GetPricesAsFormattedString(ctx, symbols, false)

	// Best Effort.
	s.ChannelMessageSend(m.ChannelID, pricesMsg)

	return nil
}

func allPriceCommand(ctx context.Context, tokens []string, s *discordgo.Session, m *discordgo.MessageCreate) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	symbols := coins.List()
	pricesMsg := priceBotSvc.GetPricesAsFormattedString(ctx, symbols, true)

	// Best Effort
	s.ChannelMessageSend(discordproto.DiscordSatoshiPriceBotChannel, pricesMsg)

	return nil
}
