package consumers

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sort"
	"strings"
	"swallowtail/libraries/util"
	coingecko "swallowtail/s.coingecko/clients"
	"swallowtail/s.twitter/clients"
	"sync"
	"syscall"
	"time"

	"github.com/monzo/slog"
)

var (
	defaultPriceBotInterval = time.Duration(1 * time.Hour)
	priceBotSymbols         = []string{
		"BTC",
		"ETH",
		"ROPE",
		"LTC",
		"OCEAN",
		"RSR",
		"NOIA",
		"HTR",
		"SOL",
		"AKT",
		"BNB",
		"ALPHA",
		"WOO",
		"ALGO",
		"AAVE",
		"RUNE",
		"SAND",
		"FET",
		"FTT",
		"RAY",
		"API3",
		"UNI",
		"1INCH",
		"BAND",
		"BAL",
		"CAKE",
		"SRM",
		"ORK",
		"AKRO",
		"SC",
		"TVK",
		"IOST",
		"BOSON",
		"FIDA",
		"OXY",
		"YFI",
		"MIR",
		"CRWNY",
		"STEP",
		"LINK",
	}
	priceBotMtx sync.RWMutex
)

func NewPriceBot(ctx context.Context, discordChannel string, coingeckoClient *coingecko.CoinGeckoClient, discordClient *clients.DiscordClient) *PriceBot {
	d := make(chan struct{}, 1)
	pb := &PriceBot{
		interval:       defaultPriceBotInterval,
		discordChannel: discordChannel,
		cgc:            coingeckoClient,
		dc:             discordClient,
		done:           d,
	}
	go func() {
		defer pb.Done(ctx)
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
		select {
		case <-sc:
			return
		case <-d:
			return
		}
	}()
	return pb
}

type PriceBot struct {
	interval       time.Duration
	discordChannel string
	cgc            *coingecko.CoinGeckoClient
	dc             *clients.DiscordClient
	done           chan struct{}
}

type PriceBotPrice struct {
	Price  float64
	Symbol string
}

type priceBotPriceList []*PriceBotPrice

func (p priceBotPriceList) Len() int {
	return len(p)
}

func (p priceBotPriceList) Less(i, j int) bool {
	return p[i].Symbol < p[j].Symbol
}

func (p priceBotPriceList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (pb *PriceBot) Start(ctx context.Context) {
	slog.Info(ctx, "Starting Price bot on a %v schedule", defaultPriceBotInterval)
	t := time.NewTicker(defaultPriceBotInterval)
	for {
		select {
		case <-t.C:
			symbols := getPriceBotSymbols()
			var (
				// Move to chan rather than mutex.
				prices = []*PriceBotPrice{}
				mtx    sync.Mutex
				wg     sync.WaitGroup
			)
			slog.Trace(ctx, "Price bot gathering prices")
			for _, symbol := range symbols {
				symbol := symbol
				go func() {
					wg.Add(1)
					defer wg.Done()
					price, err := pb.cgc.GetCurrentPriceFromSymbol(ctx, symbol, "usd")
					if err != nil {
						slog.Error(ctx, "Failed to retrieve price", map[string]string{
							"symbol":    symbol,
							"error_msg": err.Error(),
						})
						return
					}
					slog.Info(ctx, "Price bot received price", map[string]string{
						"symbol": symbol,
					})
					mtx.Lock()
					prices = append(prices, &PriceBotPrice{
						Price:  price,
						Symbol: symbol,
					})
					mtx.Unlock()
				}()
			}
			wg.Wait()
			msg := buildDiscordMessageFromPrices(prices)
			if msg == "" {
				slog.Warn(ctx, "Price bot prices empty")
				continue
			}
			err := pb.dc.PostToChannel(ctx, pb.discordChannel, msg)
			if err != nil {
				slog.Error(ctx, "Failed to post to discord", map[string]string{
					"discord_channel": pb.discordChannel,
				})
				continue
			}
			slog.Trace(ctx, "Price bot posted to discord.", nil)

		case <-pb.done:
			return
		}
	}
}

func (pb *PriceBot) Done(ctx context.Context) {
	defer slog.Info(ctx, "Cancelling price bot")
	select {
	case pb.done <- struct{}{}:
		return
	}
}

func getPriceBotSymbols() []string {
	priceBotMtx.RLock()
	defer priceBotMtx.RUnlock()
	return priceBotSymbols
}

func buildDiscordMessageFromPrices(prices []*PriceBotPrice) string {
	if len(prices) == 0 {
		return ""
	}

	sort.Sort(priceBotPriceList(prices))

	lines := []string{}
	for _, price := range prices {
		fp, err := util.FormatPriceAsString(price.Price)
		if err != nil {
			fp = "N/A"
		}
		lines = append(lines, fmt.Sprintf("[%s]: %s", price.Symbol, fp))
	}

	base := strings.Join(lines, "\n")
	greeting := fmt.Sprintf("\n:robot: **Price bot hourly update** :robot:\n[%v]\n\nPlease ping **@ajperkins** if you'd like a coin adding.\n", time.Now())
	return fmt.Sprintf("%s%s", greeting, monospaceWrapper(base))
}

func monospaceWrapper(s string) string {
	return fmt.Sprintf("```%s```", s)
}
