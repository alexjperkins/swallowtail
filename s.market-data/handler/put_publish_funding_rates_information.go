package handler

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"swallowtail/libraries/gerrors"
	accountproto "swallowtail/s.account/proto"
	discordproto "swallowtail/s.discord/proto"
	"swallowtail/s.market-data/assets"
	marketdataproto "swallowtail/s.market-data/proto"

	"github.com/monzo/slog"
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
	Exchange    accountproto.ExchangeType
	Symbol      string
	FundingRate float64
}

// PublishFundingRateInformation ...
func (s *MarketDataService) PublishFundingRateInformation(
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
			}

			fundingRate, err := handler(ctx, asset.Symbol)
			if err != nil {
				slog.Error(ctx, "Failed to get funding rate from: %v for %s", asset.Exchange, asset.Symbol)
				return
			}

			mu.Lock()
			defer mu.Unlock()
			fundingRates = append(fundingRates, &FundingRateInfo{
				Exchange:    asset.Exchange,
				Symbol:      asset.Symbol,
				FundingRate: fundingRate,
			})
		}()
	}

	wg.Wait()

	sort.Slice(fundingRates, func(i, j int) bool {
		if fundingRates[i].Symbol < fundingRates[j].Symbol {
			return true
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
	sb.WriteString(fmt.Sprintf(":robot:    Hourly Funding Rates [%v]    :orangutan:", now))

	var prevSymbol string
	for _, fr := range fundingRates {
		switch {
		case prevSymbol == "":
			prevSymbol = fr.Symbol
		case prevSymbol != fr.Symbol:
			sb.WriteString("\n")
			prevSymbol = fr.Symbol
		}

		var emoji string
		switch {
		case fr.FundingRate > 0:
			emoji = ":red_circle:"
		case fr.FundingRate < 0:
			emoji = ":green_circle:"
		default:
			emoji = ":orange_circle:"
		}

		sb.WriteString(
			fmt.Sprintf(
				"%s %s%s: %s %s %.3f\n",
				emoji,
				strings.ToTitle(fr.Exchange.String()),
				strings.Repeat(" ", exchangeIndent-len(fr.Exchange.String())),
				fr.Symbol,
				strings.Repeat(" ", symbolsIndent-len(fr.Symbol)),
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
