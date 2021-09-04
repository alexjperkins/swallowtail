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
	priceCommandUsage = `!price <[symbols... | all ]>`
)

var (
	priceBotSvc pricebot.PriceBotService
)

func init() {
	register(priceCommandID, &Command{
		ID:          priceCommandID,
		Usage:       priceCommandUsage,
		Handler:     priceCommand,
		Description: "Fetches the latest price from coingecko for the symbols provided. Pass `all` to republish the pricebot.",
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
	_, err := s.ChannelMessageSend(m.ChannelID, pricesMsg)
	return err
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
