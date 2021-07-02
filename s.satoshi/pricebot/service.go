package pricebot

import (
	"context"
	"sort"
	coingecko "swallowtail/s.coingecko/clients"
	"sync"

	"github.com/monzo/slog"
	"github.com/monzo/terrors"
)

var (
	// To aid with mocking in tests.
	defaultCoingeckoClient = coingecko.New
)

type PriceBotService interface {
	GetPrice(ctx context.Context, symbol string) (*PriceBotPrice, error)
	GetPrices(ctx context.Context, symbols []string) []*PriceBotPrice
	GetPricesAsFormattedString(ctx context.Context, symbols []string, withGreeting bool) string
}

type PriceBotPrice struct {
	Price  float64
	Symbol string
}

func NewService(ctx context.Context) PriceBotService {
	pb := &priceBotService{
		cgc: defaultCoingeckoClient(ctx),
	}
	return pb
}

type priceBotService struct {
	cgc coingecko.CoinGeckoClient
}

func (p *priceBotService) GetPricesAsFormattedString(ctx context.Context, symbols []string, withGreeting bool) string {
	prices := p.GetPrices(ctx, symbols)
	return buildMessage(prices, withGreeting)
}

func (p *priceBotService) GetPrices(ctx context.Context, symbols []string) []*PriceBotPrice {
	var (
		prices = make([]*PriceBotPrice, len(symbols))
		wg     sync.WaitGroup
	)
	sort.Strings(symbols)
	for i, symbol := range symbols {
		i, symbol := i, symbol
		wg.Add(1)
		go func() {
			defer wg.Done()
			price, err := p.GetPrice(ctx, symbol)
			if err != nil {
				// Best effort
				slog.Info(ctx, "Pricebot failed to retreive price for %s; [%v]", symbol, err)
			}
			prices[i] = price
		}()
	}
	wg.Wait()
	return prices
}

func (p *priceBotService) GetPrice(ctx context.Context, symbol string) (*PriceBotPrice, error) {
	price, err := p.cgc.GetCurrentPriceFromSymbol(ctx, symbol, "usd")
	if err != nil {
		return &PriceBotPrice{
				Symbol: symbol,
			}, terrors.Augment(err, "Pricebot failed to retreive price", map[string]string{
				"symbol": symbol,
			})
	}
	slog.Info(ctx, "Price bot received price", map[string]string{
		"symbol": symbol,
	})
	return &PriceBotPrice{
		Price:  price,
		Symbol: symbol,
	}, nil
}
