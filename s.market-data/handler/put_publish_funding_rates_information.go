package handler

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/monzo/slog"

	"swallowtail/libraries/gerrors"
	accountproto "swallowtail/s.account/proto"
	discordproto "swallowtail/s.discord/proto"
	"swallowtail/s.market-data/assets"
	marketdataproto "swallowtail/s.market-data/proto"
)

var (
	fundingRatesAssets = assets.FundingRateAssets
)

var (
	exchanges = []accountproto.ExchangeType{
		accountproto.ExchangeType_BINANCE,
		accountproto.ExchangeType_FTX,
	}
)

// FundingRateInfo ...
type FundingRateInfo struct {
	Exchange        accountproto.ExchangeType
	Symbol          string
	HumanizedSymbol string
	FundingRate     float64
}

// PublishFundingRatesInformation ...
func (s *MarketDataService) PublishFundingRatesInformation(
	ctx context.Context, in *marketdataproto.PublishFundingRatesInformationRequest,
) (*marketdataproto.PublishFundingRatesInformationResponse, error) {
	slog.Trace(ctx, "Market data publishing funding rates information")

	var (
		fundingRates = make([]*FundingRateInfo, 0, len(fundingRatesAssets))
		wg           sync.WaitGroup
		mu           sync.RWMutex
	)
	for _, asset := range fundingRatesAssets {
		asset := asset

		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(jitter(0, 59))

			var handler func(ctx context.Context, symbol string) (float64, error)
			switch asset.Exchange {
			case accountproto.ExchangeType_BINANCE:
				handler = getFundingRateFromBinance
			case accountproto.ExchangeType_FTX:
				handler = getFundingRateFromFTX
			case accountproto.ExchangeType_BITFINEX:
				handler = getFundingRateFromBitfinex
			}

			fundingRate, err := handler(ctx, asset.Symbol)
			if err != nil {
				slog.Error(ctx, "Failed to get funding rate from: %v for %s: %v", asset.Exchange, asset.Symbol, err)
				return
			}

			mu.Lock()
			defer mu.Unlock()
			fundingRates = append(fundingRates, &FundingRateInfo{
				Exchange:        asset.Exchange,
				Symbol:          asset.Symbol,
				HumanizedSymbol: asset.HumanizedSymbol,
				FundingRate:     fundingRate * 100,
			})
		}()
	}

	wg.Wait()

	sort.Slice(fundingRates, func(i, j int) bool {
		if fundingRates[i].Symbol < fundingRates[j].Symbol {
			return true
		}
		if fundingRates[i].Symbol > fundingRates[j].Symbol {
			return false
		}

		return fundingRates[i].Exchange < fundingRates[j].Exchange
	})

	var exchangeIndent int
	for _, ex := range exchanges {
		if len(ex.String()) > exchangeIndent {
			exchangeIndent = len(ex.String())
		}
	}

	var symbolsIndent int
	for _, fr := range fundingRates {
		if len(fr.Symbol) > symbolsIndent {
			symbolsIndent = len(fr.Symbol)
		}
	}

	now := time.Now().UTC().Truncate(time.Hour)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(":robot:    `Market Data: Hourly Funding Rates [%v]`    :orangutan:\n", now))

	for _, fr := range fundingRates {
		var (
			exchangeInfo = assets.GetFundingRateCoefficientByExchange(fr.Exchange)
			emoji        string
		)
		switch {
		case fr.FundingRate > exchangeInfo.HigherBound:
			emoji = ":red_circle:"
		case fr.FundingRate < exchangeInfo.LowerBound:
			emoji = ":green_circle:"
		default:
			emoji = ":orange_circle:"
		}

		symbol := fr.Symbol
		if fr.HumanizedSymbol != "" {
			symbol = fr.HumanizedSymbol
		}

		sb.WriteString(
			fmt.Sprintf(
				"\n%s `[%s]:    %s %s %s %.4f`",
				emoji,
				fr.Symbol,
				strings.Repeat(" ", symbolsIndent-len(symbol)),
				strings.ToTitle(fr.Exchange.String()),
				strings.Repeat(" ", exchangeIndent-len(fr.Exchange.String())),
				fr.FundingRate,
			),
		)
	}

	idempotencyKey := fmt.Sprintf("fundingrate-%v", now)
	if err := publishToDiscord(ctx, sb.String(), discordproto.DiscordSatoshiPriceBotChannel, idempotencyKey); err != nil {
		return nil, gerrors.Augment(err, "failed_to_publish_funding_rate_information", map[string]string{
			"idempotency_key": idempotencyKey,
		})
	}

	return &marketdataproto.PublishFundingRatesInformationResponse{}, nil
}
