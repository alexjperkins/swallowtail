package consumers

import (
	"context"
	"fmt"
	"math/rand"
	discordproto "swallowtail/s.discord/proto"
	"swallowtail/s.satoshi/coins"
	"time"

	"github.com/monzo/slog"
)

const (
	volatilityConsumerID = "volatility-consumer"
	volatilityInterval   = time.Duration(5 * time.Minute)

	// Jitter defaults
	volatilityJitterMax      = 3
	volatilityJitterDuration = time.Minute

	// Config
	highVolatilityDelta = 0.03
	midVolatilityDelta  = 0.02
	lowVolatilityDelta  = 0.0075
)

var (
	defaultVolatilitySymbols []string
	lowVolatilitySet         = map[string]bool{
		"BTC": true,
	}
	highVolatilitySet = map[string]bool{
		"ALPHA": true,
	}
)

type priceAction struct {
	curr float64
	prev float64
}

func init() {
	defaultVolatilitySymbols = coins.List()
	register(volatilityConsumerID, VolatilityConumser{
		Active: true,
	})
}

type VolatilityConumser struct {
	Active bool
}

func (v VolatilityConumser) Receiver(ctx context.Context, c chan *ConsumerMessage, d chan struct{}, withJitter bool) {
	for _, symbol := range defaultVolatilitySymbols {
		symbol := symbol
		go func() {
			if withJitter {
				// Sleep for jitter amount
				time.Sleep(time.Duration(rand.Intn(volatilityJitterMax)) * volatilityJitterDuration)
			}
			var (
				pa = priceAction{}
				t  = time.NewTicker(volatilityInterval)
			)
			defer slog.Warn(ctx, "Volatility consumer stop signal received.")
			for {
				select {
				case <-t.C:
					curr, err := getAssetLatestPrice(ctx, symbol, "usd")
					if err != nil {
						// Best effort
						slog.Error(ctx, "Failed to get current price: %v", err)
						continue
					}
					pa.curr = curr

					difference := pa.percentageDiff()
					trigger := getVolatilityDelta(symbol)
					switch {
					case difference > trigger:
						publishVolatilityMsg(c, symbol, pa, difference, true, v.IsActive())
					case difference < trigger:
						publishVolatilityMsg(c, symbol, pa, difference, false, v.IsActive())
					default:
						continue
					}
					pa.prev = pa.curr

				case <-d:
					return
				case <-ctx.Done():
					return
				}
			}
		}()
	}
}

func (v VolatilityConumser) IsActive() bool {
	return v.Active
}

func publishVolatilityMsg(c chan *ConsumerMessage, symbol string, pa priceAction, difference float64, increasing, isActive bool) {
	direction := "MOON"
	incOrDecr := "INCREASE"
	if !increasing {
		direction = "REKT"
		incOrDecr = "DECREASE"
	}
	msg := fmt.Sprintf(
		":rotating_light: [%s] %s %.3f%% %s in %v from %.4f :arrow_forward: %.4f",
		symbol,
		direction,
		abs(difference)*100,
		incOrDecr,
		volatilityInterval,
		pa.prev,
		pa.curr,
	)
	select {
	case c <- &ConsumerMessage{
		ConsumerID:       volatilityConsumerID,
		DiscordChannelID: discordproto.DiscordSatoshiAlertsChannel,
		Message:          msg,
		IsActive:         isActive,
		Created:          time.Now(),
	}:
	}
}

func (pa *priceAction) percentageDiff() float64 {
	if pa.prev == 0.0 || pa.curr == 0.0 {
		return 0.0
	}

	return (pa.curr / pa.prev) - 1
}

func getVolatilityDelta(symbol string) float64 {
	if _, ok := highVolatilitySet[symbol]; ok {
		return highVolatilityDelta
	}

	if _, ok := lowVolatilitySet[symbol]; ok {
		return lowVolatilityDelta
	}

	return midVolatilityDelta
}

func abs(v float64) float64 {
	if v < 0 {
		return -v
	}

	return v
}
