package consumers

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	coingecko "swallowtail/s.coingecko/clients"
	"swallowtail/s.satoshi/clients"
	"sync"
	"syscall"
	"time"

	"github.com/monzo/slog"
)

var (
	defaultATHJitterDuration     = time.Minute
	defaultATHJitterMax          = 3
	defaultApproachingATHTrigger = 0.02

	coinIDMtx              sync.Mutex
	defaultATHAlertCoinIDs = map[string]*CoinInfo{
		"1INCH": {
			ID: "1inch",
		},
		"AAVE": {
			ID: "aave",
		},
		"ALGO": {
			ID: "algorand",
		},
		"ALPHA": {
			ID: "alpha-finance",
		},
		"BAND": {
			ID: "band-protocol",
		},
		"BNB": {
			ID: "binancecoin",
		},
		"BTC": {
			ID: "bitcoin",
		},
		"CAKE": {
			ID: "pancakeswap-token",
		},
		"DOT": {
			ID: "polkadot",
		},
		"ETH": {
			ID: "ethereum",
		},
		"LINK": {
			ID: "chainlink",
		},
		"LTC": {
			ID: "litecoin",
		},
		"OCEAN": {
			ID: "ocean-protocol",
		},
		"RSR": {
			ID: "reserve-rights-token",
		},
		"SOL": {
			ID: "solana",
		},
		"SRM": {
			ID: "serum",
		},
		"UNI": {
			ID: "uniswap",
		},
	}
)

type CoinInfo struct {
	ID string
}

func NewATHAlerter(symbol string, interval time.Duration, discordClient *clients.DiscordClient, coingeckoClient *coingecko.CoinGeckoClient, withJitter bool) *ATHAlerter {
	done := make(chan struct{}, 1)
	errCh := make(chan error, 32)
	a := &ATHAlerter{
		symbol:     symbol,
		interval:   interval,
		cgc:        coingeckoClient,
		dc:         discordClient,
		done:       done,
		errCh:      errCh,
		withJitter: withJitter,
	}

	go func() {
		defer slog.Info(nil, "Closing down ATH alerter for: %s", symbol)
		defer a.Done()
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
		select {
		case <-sc:
			return
		case <-done:
			return
		}
	}()
	return a
}

type ATHAlerter struct {
	symbol     string
	interval   time.Duration
	cgc        *coingecko.CoinGeckoClient
	dc         *clients.DiscordClient
	done       chan struct{}
	errCh      chan error
	withJitter bool
}

func (a *ATHAlerter) Run(ctx context.Context) {
	var (
		currentPrice float64
		currentATH   float64
	)
	id, ok := GetCoinInfoFromSymbol(a.symbol)
	if !ok {
		slog.Info(context.TODO(), "Metadata not stored for %s, cannot query `coingecko`", a.symbol)
		a.Done()
		return
	}
	slog.Info(context.TODO(), "Starting ATH Alerter for %s every %v", a.symbol, a.interval)

	currentATH, err := a.cgc.GetATHFromID(ctx, id)
	if err != nil {
		slog.Error(context.TODO(), "Failed to get ATH for %s, error: %s", a.symbol, err)
		return
	}

	slog.Info(context.TODO(), "Received ATH for %s: %.4f", a.symbol, currentATH)

	t := time.NewTicker(a.interval)
	// Sleep for random time
	if a.withJitter {
		time.Sleep(time.Duration(rand.Intn(defaultATHJitterMax)) * defaultATHJitterDuration)
	}
	for {
		select {
		case <-t.C:
			currentPrice, err = a.cgc.GetCurrentPriceFromID(ctx, id, "usd")
			slog.Info(context.TODO(), "Recieved current price for: %s, %.4f", a.symbol, currentPrice)
			if err != nil {
				slog.Info(context.TODO(), "Failed to get the current price %s", err.Error())
				continue
			}

			// Check if approaching
			if isApproachingATH(currentPrice, currentATH) {
				a.dc.PostToChannel(context.TODO(), clients.DiscordAlertsChannel, fmt.Sprintf(":rotating_light: ATH Alert: %s is approaching ATH of %.4f in 3, 2, 1...", a.symbol, currentATH))
				slog.Info(context.TODO(), ":rotating_light: ATH Alert: %s is approaching ATH of %.4f in 3, 2, 1...", a.symbol, currentATH)
				continue
			}
			if currentPrice < currentATH {
				continue
			}

			slog.Info(context.TODO(), "%s ATH alert triggered, previous: %.4f new: %.4f", a.symbol, currentATH, currentPrice)
			// Best effort
			a.dc.PostToChannel(context.TODO(), clients.DiscordAlertsChannel, fmt.Sprintf(":rocket: New ATH alert: %s, previous %.4f, new %.4f :new_moon:", a.symbol, currentATH, currentPrice))

			currentATH = currentPrice
		case <-a.done:
			return
		}
	}
}

func (a *ATHAlerter) Done() {
	defer slog.Info(context.TODO(), "Cancelling ATH for %s", a.symbol)
	select {
	case a.done <- struct{}{}:
		return
	}
}

func GetDefaultATHAlertCoins() map[string]*CoinInfo {
	coinIDMtx.Lock()
	defer coinIDMtx.Unlock()

	nm := map[string]*CoinInfo{}

	for k, v := range defaultATHAlertCoinIDs {
		nm[k] = v
	}

	return nm
}

func GetCoinInfoFromSymbol(symbol string) (string, bool) {
	coinIDMtx.Lock()
	defer coinIDMtx.Unlock()

	ci, ok := defaultATHAlertCoinIDs[strings.ToUpper(symbol)]
	return ci.ID, ok
}

func isApproachingATH(p, ath float64) bool {
	distance := p/ath - 1
	if distance > 0 && distance < defaultApproachingATHTrigger {
		return true
	}
	return false
}
