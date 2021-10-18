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

func priceCommand(ctx context.Context, tokens []string, s *discordgo.Session, m *discordgo.MessageCreate) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	symbols := tokens

	var (
		cache = make(map[string]*coingeckoproto.GetAssetLatestPriceBySymbolResponse)
		wg    sync.WaitGroup
		mu    sync.Mutex
	)
	for _, symbol := range symbols {
		symbol := symbol
		wg.Add(1)
		go func() {
			defer wg.Done()

			rsp, err := (&coingeckoproto.GetAssetLatestPriceBySymbolRequest{
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

			mu.Lock()
			defer mu.Unlock()
			cache[symbol] = rsp
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

		sb.WriteString(fmt.Sprintf("%s `[%s] %.3f USDT 24h: %.2f%%\n`", emoji, k, v.LatestPrice, v.PercentagePriceChange_24H))
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
