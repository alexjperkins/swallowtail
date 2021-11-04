package commands

import (
	"context"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"sync"
	"time"

	"swallowtail/libraries/gerrors"
	binanceproto "swallowtail/s.binance/proto"
	coingeckoproto "swallowtail/s.coingecko/proto"

	"github.com/bwmarrin/discordgo"
	"github.com/monzo/slog"
)

const (
	priceCommandID    = "price"
	priceCommandUsage = `!price <[symbols, ...]>`
)

func init() {
	register(priceCommandID, &Command{
		ID:          priceCommandID,
		Usage:       priceCommandUsage,
		Handler:     priceCommand,
		Description: "Fetches the latest price from coingecko for the symbols provided.",
	})
}

type InstrumentInfo struct {
	CurrentPrice              float64
	FundingRate               float64
	Symbol                    string
	AssetPair                 string
	PercentagePriceChange_24H float64
}

func priceCommand(ctx context.Context, tokens []string, s *discordgo.Session, m *discordgo.MessageCreate) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	symbols := tokens

	var (
		cache = make(map[string]*InstrumentInfo)
		wg    sync.WaitGroup
		mu    sync.Mutex
	)
	for _, symbol := range symbols {
		symbol := symbol
		wg.Add(1)
		go func() {
			defer wg.Done()

			cgRsp, err := (&coingeckoproto.GetAssetLatestPriceBySymbolRequest{
				AssetSymbol: symbol,
				AssetPair:   "usd",
			}).SendWithTimeout(ctx, 1*time.Minute).Response()
			if err != nil {
				slog.Warn(ctx, "Failed to fetch coingecko price for price command", map[string]string{
					"symbol": symbol,
					"error":  err.Error(),
				})
				return
			}

			// Parse funding rate if we can. Best effort.
			var fundingRate float64
			rsp, _ := (&binanceproto.GetFundingRatesRequest{
				Symbol: fmt.Sprintf("%sUSDT", strings.ToUpper(symbol)),
				Limit:  1,
			}).SendWithTimeout(ctx, 1*time.Minute).Response()
			if rsp != nil && len(rsp.GetFundingRates()) > 0 {
				fundingRate = float64(rsp.GetFundingRates()[0].FundingRate)
			}

			mu.Lock()
			defer mu.Unlock()
			cache[symbol] = &InstrumentInfo{
				CurrentPrice:              float64(cgRsp.LatestPrice),
				PercentagePriceChange_24H: float64(cgRsp.PercentagePriceChange_24H),
				FundingRate:               fundingRate,
			}
		}()
	}

	wg.Wait()

	if len(cache) == 0 {
		// Best Effort
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":wave: <@%s> Sorry, i wasn't able to get price info for any symbols [%s]:disappointed:", m.Author.ID, strings.Join(symbols, ",")))
		return nil
	}

	var keys = make([]string, 0, len(cache))
	for k := range cache {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		v, ok := cache[k]
		if !ok {
			// We're being overly defensive here.
			continue
		}

		var emoji string
		switch {
		case v.PercentagePriceChange_24H > 0:
			emoji = ":green_square:"
		case v.PercentagePriceChange_24H < 0:
			emoji = ":red_square:"
		default:
			emoji = ":black_large_square:"
		}

		sb.WriteString(fmt.Sprintf("%s `[%s] %.3f USDT 24h: %.2f%%  Funding Rate: %.4f%%\n`", emoji, k, v.CurrentPrice, v.PercentagePriceChange_24H, v.FundingRate*100))
	}

	if _, err := s.ChannelMessageSend(m.ChannelID, sb.String()); err != nil {
		return gerrors.Augment(err, "failed_to_execute_price_command", nil)
	}

	return nil
}

func jitter() time.Duration {
	rand.Seed(time.Now().UnixNano())
	return time.Duration(rand.Intn(5)) * time.Second
}
