package consumers

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"swallowtail/libraries/util"
	"swallowtail/s.satoshi/clients"
	"syscall"
	"time"

	"github.com/monzo/slog"
)

var (
	percentageIncreaseTrigger = 0.015
	percetangeDecreaseTrigger = -0.015

	highVolatilityTrigger   = 0.03
	mediumVolatilityTrigger = 0.02
	lowVolatilityTrigger    = 0.0075

	defaultVolatilityTrigger = 0.02

	increasingGIF = "https://tenor.com/view/turbo-time-jinglealltheway-arnold-gif-11070880"
	decreasingGIF = "https://tenor.com/view/stop-hes-dead-already-gif-11313631"

	volatiliyMapping = map[string]float64{
		"BTCUSDT":   lowVolatilityTrigger,
		"LINKUSDT":  mediumVolatilityTrigger,
		"ALPHAUSDT": highVolatilityTrigger,
	}

	jitterDefaultInterval = 3
	jitterDefaultDuration = time.Minute
)

func NewVolatilityAlerter(symbol string, binanceClient *clients.BinanceClient, discordClient *clients.DiscordClient, channel string, interval time.Duration, withJitter bool) *VolatilityAlerter {
	done := make(chan struct{}, 1)
	errCh := make(chan error, 32)

	a := &VolatilityAlerter{
		symbol:     symbol,
		bc:         binanceClient,
		dc:         discordClient,
		interval:   interval,
		withJitter: withJitter,
		channel:    channel,
		done:       done,
		errCh:      errCh,
	}

	go func() {
		defer slog.Info(nil, "Closing down discord ws")
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

type VolatilityAlerter struct {
	symbol     string
	bc         *clients.BinanceClient
	dc         *clients.DiscordClient
	interval   time.Duration
	withJitter bool
	channel    string
	done       chan struct{}
	errCh      chan error
}

type priceAction struct {
	curr float64
	prev float64
}

func (pa *priceAction) percentageDiff() float64 {
	if pa.prev == 0.0 || pa.curr == 0.0 {
		return 0.0
	}
	return (pa.curr / pa.prev) - 1

}

func (a *VolatilityAlerter) Run(ctx context.Context) {
	if a.withJitter {
		// Sleep for jitter amount
		time.Sleep(time.Duration(rand.Intn(jitterDefaultInterval)) * jitterDefaultDuration)
	}

	slog.Info(ctx, "Starting alerter for %s with interval %v, jitter set: %s", a.symbol, a.interval, a.withJitter)
	t := time.NewTicker(a.interval)
	defer slog.Info(ctx, "Closing down BTC Alerter.")

	var pa = priceAction{}
	for {
		select {
		case <-t.C:
			slog.Info(ctx, "Fetching price from Binance")
			latestValue, err := a.bc.GetPrice(ctx, a.symbol)
			if err != nil {
				select {
				case a.errCh <- err:
					continue
				}
			}

			strCurrentPrice := latestValue.MarkPrice
			f, err := strconv.ParseFloat(strCurrentPrice, 64)
			if err != nil {
				select {
				case a.errCh <- err:
					continue
				}
			}
			pa.curr = f

			diff := pa.percentageDiff()
			trigger := abs(getTrigger(a.symbol))
			switch {
			case diff < trigger && diff > -trigger:
				slog.Info(ctx, "%s: low volatiliy: %f, %f -> %f", a.symbol, diff, pa.prev, pa.curr)
			case diff > trigger:
				slog.Info(ctx, "%s: Price Alert %v, %v", a.symbol, diff, trigger)
				a.dc.PostToChannel(ctx, a.channel, percentageIncreaseMsg(a.symbol, pa.prev, pa.curr, diff, a.interval))
			case diff < trigger:
				slog.Info(ctx, "%s: Price Alert %v, %v", a.symbol, diff, trigger)
				a.dc.PostToChannel(ctx, a.channel, percentageDecreaseMsg(a.symbol, pa.prev, pa.curr, diff, a.interval))
			}

			pa.prev = pa.curr
			continue

		case <-a.done:
			return
		case <-ctx.Done():
			return
		}

	}
}

func (a *VolatilityAlerter) Done() {
	slog.Info(context.TODO(), "Alerter done signal recieved.")
	select {
	case a.done <- struct{}{}:
		return
	}
}

func percentageIncreaseMsg(symbol string, prev, current, diff float64, interval time.Duration) string {
	c, err := util.FormatPriceAsString(current)
	if err != nil {
		c = ""
	}
	p, err := util.FormatPriceAsString(prev)
	if err != nil {
		p = ""
	}
	return fmt.Sprintf(
		":new_moon_with_face: %s **MOON** %.3f%% INCREASE in %v from %s :arrow_forward: %s \n %s",
		symbol,
		abs(diff*100),
		interval,
		p,
		c,
		increasingGIF,
	)
}

func percentageDecreaseMsg(symbol string, prev, current, diff float64, interval time.Duration) string {
	c, err := util.FormatPriceAsString(current)
	if err != nil {
		c = ""
	}
	p, err := util.FormatPriceAsString(prev)
	if err != nil {
		p = ""
	}
	return fmt.Sprintf(
		":warning: %s **REKT** %.3f%% DECREASE in %v from %s :arrow_forward: %s \n %s",
		symbol,
		abs(diff*100),
		interval,
		p,
		c,
		decreasingGIF,
	)
}

func abs(f float64) float64 {
	if f < 0.0 {
		return f * -1
	}
	return f
}

func getTrigger(symbol string) float64 {
	if v, ok := volatiliyMapping[symbol]; ok {
		return v
	}
	return defaultVolatilityTrigger
}
