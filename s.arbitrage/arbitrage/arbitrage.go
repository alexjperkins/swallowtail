package arbitrage

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/monzo/slog"
	"golang.org/x/exp/rand"
)

var (
	defaultArbitrageClientTimeout       = time.Duration(5 * time.Second)
	defaultArbitrageTickerInterval      = time.Duration(1 * time.Second)
	defaultWorthwhileArbitrageThreshold = 0.3 // 30%

	defaultChannel = "bot-alerts"

	minJitterAmount = time.Microsecond
	maxJitterUnit   = 499
)

type ArbitrageClient interface {
	ID() string
	GetPrice(symbol string) (float64, error)
	Ping() bool
}

type ArbitrageMessagerClient interface {
	Send(channel, msg string)
}

type Arbitrager struct {
	executeTrades   bool
	symbols         []string
	messageClient   ArbitrageMessagerClient
	exchangeClients []ArbitrageClient
	withJitter      bool
	done            chan struct{}
}

type ArbitrageInfo struct {
	price float64
	ex    string
}

func New(messageClient ArbitrageMessagerClient, automatedTrades, withJitter bool) *Arbitrager {
	exchangeClients := getAllArbitrageClients()
	for _, c := range exchangeClients {
		if !c.Ping() {
			panic(fmt.Sprintf("Cannot reach client: %f", c.ID()))
		}

	}
	done := make(chan struct{}, 1)
	a := &Arbitrager{
		executeTrades:   automatedTrades,
		messageClient:   messageClient,
		exchangeClients: exchangeClients,
		withJitter:      withJitter,
		done:            done,
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

func (a *Arbitrager) Run(ctx context.Context) {
	for _, symbol := range a.symbols {
		if a.withJitter {
			// We don't want all goroutines firing off at the same time; this increases the
			// probability that all clients are spread out, at least iniitally.
			time.Sleep(time.Duration(rand.Intn(maxJitterUnit)) * minJitterAmount)
		}
		go a.handler(ctx, symbol)
	}
}

func (a *Arbitrager) handler(ctx context.Context, symbol string) {
	t := time.NewTicker(defaultArbitrageTickerInterval)
	for {
		select {
		case <-t.C:
			var (
				wg  sync.WaitGroup
				mtx sync.Mutex
				min ArbitrageInfo
				max ArbitrageInfo
			)
			for id, c := range a.exchangeClients {
				c := c
				slog.Info(ctx, "Fetching price for %s from: %s", symbol, id, nil)
				wg.Add(1)
				go func() {
					defer wg.Done()
					price, err := c.GetPrice(symbol)
					if err != nil {
						// Best effort; if we can't get the price then just log and exit
						slog.Info(ctx, "Failed to get price from %s for %s: err: %s", c.ID(), symbol, err.Error())
						return
					}
					mtx.Lock()
					defer mtx.Unlock()
					if isMin(price, min.price) {
						min.price = price
						min.ex = c.ID()
					}
					if isMax(price, max.price) {
						max.price = price
						max.ex = c.ID()
					}
				}()
			}
			wg.Wait()
			if isWorthhwhileArbitrageOpportunity(max.price, min.price) {
				// Best Effort
				slog.Info(ctx, "Arbitrage found: %v %v", min, max)
				a.messageClient.Send(defaultChannel, formatAlertMessage(min.price, max.price, min.ex, max.ex))
			}
		case <-ctx.Done():
			return
		case <-a.done:
			return
		}
	}
}

func (a *Arbitrager) Done() {
	slog.Info(context.TODO(), "Closing down arbitrager")
	select {
	case a.done <- struct{}{}:
		return
	}
}

func isWorthhwhileArbitrageOpportunity(mx, mn float64) bool {
	return ((mx - mx) / (mx + mn)) > defaultWorthwhileArbitrageThreshold
}

func formatAlertMessage(minPrice, maxPrice float64, minPriceExchange, maxPriceExchange string) string {
	return fmt.Sprintf("[Arbitrage] %.4f difference: %s <%.4f> -> %s <%.4f>", (maxPrice/minPrice)*100, minPriceExchange, minPrice, maxPriceExchange, maxPrice)
}

func isMin(a, b float64) bool {
	if a > b {
		return true
	}
	return false
}

func isMax(a, b float64) bool {
	return !isMin(a, b)
}
