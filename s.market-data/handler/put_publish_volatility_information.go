package handler

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"swallowtail/libraries/gerrors"
	"swallowtail/s.market-data/assets"
	marketdataproto "swallowtail/s.market-data/proto"
)

var (
	volatilityAssets = assets.LatestPriceAssets
	volatilityOnce   sync.Once
)

// PublishVolatilityInformation ...
func (s *MarketDataService) PublishVolatilityInformation(
	ctx context.Context, in *marketdataproto.PublishVolatilityInformationRequest,
) (*marketdataproto.PublishVolatilityInformationResponse, error) {

	return nil, gerrors.ErrUnimplemented
}

func formatVolatilityContent(asset *assets.AssetPair, latestPrice, diff float64, increasing bool) string {
	var emoji = ":chart_with_upwards_trend:"
	if !increasing {
		emoji = ":chart_with_downwards_trend:"
	}

	header := fmt.Sprintf(":rotating_light:    `High Volatility Alert: %s%s` %s    :robot:", strings.ToUpper(asset.Symbol), strings.ToUpper(asset.AssetPair), emoji)
	content := `
ASSET:        %s%s
LATEST PRICE: %.3f
15M_CHANGE :  %.2f%%
`
	formattedContent := fmt.Sprintf(content, strings.ToUpper(asset.Symbol), strings.ToUpper(asset.AssetPair), latestPrice, diff*100)
	return fmt.Sprintf("%s```%s```", header, formattedContent)
}
