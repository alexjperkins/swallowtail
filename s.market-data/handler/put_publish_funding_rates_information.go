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
	discordproto "swallowtail/s.discord/proto"
	"swallowtail/s.market-data/assets"
	marketdataproto "swallowtail/s.market-data/proto"
	tradeengineproto "swallowtail/s.trade-engine/proto"
)

var (
	fundingRatesAssets = assets.FundingRateAssets
)

var (
	venues = []tradeengineproto.VENUE{
		tradeengineproto.VENUE_BINANCE,
		tradeengineproto.VENUE_FTX,
		tradeengineproto.VENUE_BITFINEX,
	}
)

// FundingRateInfo ...
type FundingRateInfo struct {
	Venue           tradeengineproto.VENUE
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
			switch asset.Venue {
			case tradeengineproto.VENUE_BINANCE:
				handler = getFundingRateFromBinance
			case tradeengineproto.VENUE_FTX:
				handler = getFundingRateFromFTX
			case tradeengineproto.VENUE_BITFINEX:
				handler = getFundingRateFromBitfinex
			}

			fundingRate, err := handler(ctx, asset.Symbol)
			if err != nil {
				slog.Error(ctx, "Failed to get funding rate from: %v for %s: %v", asset.Venue, asset.Symbol, err)
				return
			}

			mu.Lock()
			defer mu.Unlock()
			fundingRates = append(fundingRates, &FundingRateInfo{
				Venue:           asset.Venue,
				Symbol:          asset.Symbol,
				HumanizedSymbol: asset.HumanizedSymbol,
				FundingRate:     fundingRate * 100,
			})
		}()
	}

	wg.Wait()

	sort.Slice(fundingRates, func(i, j int) bool {
		fi, fj := fundingRates[i], fundingRates[j]

		var si, sj string
		switch {
		case fi.HumanizedSymbol != "":
			si = fi.HumanizedSymbol
		default:
			si = fi.Symbol
		}

		switch {
		case fj.HumanizedSymbol != "":
			sj = fj.HumanizedSymbol
		default:
			sj = fj.Symbol
		}

		if si < sj {
			return true
		}
		if si > sj {
			return false
		}

		return fi.Venue < fj.Venue
	})

	var exchangeIndent int
	for _, venue := range venues {
		if len(venue.String()) > exchangeIndent {
			exchangeIndent = len(venue.String())
		}
	}

	var symbolsIndent int
	for _, fr := range fundingRates {
		symbol := fr.Symbol
		if fr.HumanizedSymbol != "" {
			symbol = fr.HumanizedSymbol
		}

		if len(symbol) > symbolsIndent {
			symbolsIndent = len(symbol)
		}
	}

	now := time.Now().UTC().Truncate(time.Hour)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(":robot:    `Market Data: Hourly Funding Rates [%v]`    :orangutan:\n", now))

	for _, fr := range fundingRates {
		var (
			exchangeInfo = assets.GetFundingRateCoefficientByVenue(fr.Venue)
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
				symbol,
				addPadding(symbolsIndent-len(symbol)),
				strings.ToTitle(fr.Venue.String()),
				addPadding(exchangeIndent-len(fr.Venue.String())),
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
