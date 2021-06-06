package satoshi

import (
	"context"
	"sort"
	discordproto "swallowtail/s.discord/proto"
	"swallowtail/s.satoshi/coins"
	"swallowtail/s.satoshi/pricebot"
	"sync"
	"time"

	"github.com/monzo/slog"
)

var (
	PriceBotConsumerID = "price-bot-consumer"

	defaultTickInterval = time.Duration(1 * time.Hour)
	priceBotSymbols     = []string{}
	priceBotSymbolsOnce sync.Once
)

func init() {
	// registerSatoshiConsumer(PriceBotConsumerID, PriceBotConsumer{ Active: true })
	priceBotSymbolsOnce.Do(func() {
		priceBotSymbols = coins.GetCoinSymbols()
		sort.Strings(priceBotSymbols)
	})
}

// PriceBotConsumer
type PriceBotConsumer struct {
	Active bool
}

func (p PriceBotConsumer) Receiver(ctx context.Context, c chan *SatoshiConsumerMessage, d chan struct{}, _ bool) {
	svc := pricebot.NewService(ctx)
	t := time.NewTicker(defaultTickInterval)
	defer slog.Warn(ctx, "Consumer [%s] stopping", PriceBotConsumerID)
	for {
		select {
		case <-t.C:
			priceMsg := svc.GetPricesAsFormattedString(ctx, priceBotSymbols, true)
			msg := &SatoshiConsumerMessage{
				ConsumerID:       PriceBotConsumerID,
				DiscordChannelID: discordproto.DiscordSatoshiPriceBotChannel,
				Message:          priceMsg,
				Created:          time.Now(),
				IsActive:         p.IsActive(),
			}
			select {
			case c <- msg:
			default:
				slog.Warn(ctx, "Failed to publish pricebot messageo; block satoshi consumer channel")

			}
		case <-ctx.Done():
		case <-d:
		}
	}
}

func (p PriceBotConsumer) IsActive() bool {
	return p.Active
}
