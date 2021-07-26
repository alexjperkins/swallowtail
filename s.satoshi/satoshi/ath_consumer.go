package satoshi

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	coingecko "swallowtail/s.coingecko/client"
	discordproto "swallowtail/s.discord/proto"
	"swallowtail/s.satoshi/coins"

	"github.com/monzo/slog"
)

const (
	athConsumerID = "ath-consumer"

	athInterval = time.Duration(15 * time.Minute)
	// default Jitter config
	athJitterMax      = 3
	athJitterDuration = time.Minute

	approachingTriggerDelta = 0.02
)

var (
	defaultATHSymbols  []string
	athCoingeckoClient *coingecko.CoinGeckoClient
)

func init() {
	registerSatoshiConsumer(athConsumerID, ATHConsumer{
		Active: true,
	})
	defaultATHSymbols = coins.GetCoinSymbols()
}

// ATHConsumer
type ATHConsumer struct {
	Active bool
}

func (a ATHConsumer) Receiver(ctx context.Context, c chan *SatoshiConsumerMessage, d chan struct{}, withJitter bool) {
	cli := coingecko.New(ctx)

	// TODO: Move to tombstones; https://blog.labix.org/2011/10/09/death-of-goroutines-under-control
	for _, symbol := range defaultATHSymbols {
		symbol := symbol
		go func() {
			// Sleep for a random time as not to publish messages all at once.
			if withJitter {
				time.Sleep(time.Duration(rand.Intn(athJitterMax)) * athJitterDuration)
			}
			t := time.NewTicker(athInterval)
			var (
				currentPrice float64
				currentATH   float64
			)
			currentATH, err := cli.GetATHFromSymbol(ctx, symbol)
			if err != nil {
				// Best effort
				slog.Error(ctx, "Failed to fetch ATH for %s; %v", symbol, err)
			}
			defer slog.Warn(ctx, "ATH consumer stop signal received.")
			for {
				select {
				case <-t.C:
					// TODO: add directionality
					// current price is approaching our current stored ATH.
					if isApproachingATH(currentPrice, currentATH) {
						msg := fmt.Sprintf(":rotating_light: ATH Alert: %s is approaching ATH of %.4f in 3, 2, 1...", symbol, currentATH)
						publishATHMsg(c, msg, symbol, currentPrice, currentATH)
						continue
					}
					if currentPrice < currentATH {
						continue
					}
					// current price is greater than our current storedATH.
					msg := fmt.Sprintf(":rocket: New ATH alert: %s, previous %.4f, new %.4f :new_moon:", symbol, currentATH, currentPrice)
					publishATHMsg(c, msg, symbol, currentPrice, currentATH)
					currentATH = currentPrice
				case <-d:
					return
				case <-ctx.Done():
					return
				}
			}
		}()
	}
}

func (a ATHConsumer) IsActive() bool {
	return a.Active
}

func publishATHMsg(c chan<- *SatoshiConsumerMessage, msg, symbol string, currentPrice, currentATH float64) {
	created := time.Now()
	// Idempotent on the current price & current all the time & valid for an hour.
	idempotencyKey := fmt.Sprintf("%v-%v-%v", currentPrice, currentATH, created.Hour())
	select {
	case c <- &SatoshiConsumerMessage{
		ConsumerID:       athConsumerID,
		Message:          msg,
		DiscordChannelID: discordproto.DiscordSatoshiAlertsChannel,
		IdempotencyKey:   idempotencyKey,
		Created:          created,
		Metadata: map[string]string{
			"symbol":        symbol,
			"current_price": strconv.FormatFloat(currentPrice, 'f', 4, 64),
			"current_ath":   strconv.FormatFloat(currentATH, 'f', 4, 64),
		},
	}:
	default:
		slog.Warn(context.Background(), "Failed to publish ATH msg; channel blocked", map[string]string{
			"timestamp": created.String(),
			"symbol":    symbol,
		})
	}
}

func isApproachingATH(currentPrice, ath float64) bool {
	distance := (ath / currentPrice) - 1
	if distance > 0 && distance < approachingTriggerDelta {
		return true
	}
	return false
}
