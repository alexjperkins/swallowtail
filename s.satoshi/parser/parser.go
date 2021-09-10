package parser

import (
	"context"
	"strings"
	"swallowtail/libraries/gerrors"
	binanceproto "swallowtail/s.binance/proto"
	tradeengineproto "swallowtail/s.trade-engine/proto"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/monzo/slog"
)

var (
	binanceAssetPairs = map[string]bool{}
)

// TradeParser ...
type TradeParser interface {
	Parse(ctx context.Context, content string, m *discordgo.MessageCreate) (*tradeengineproto.Trade, error)
}

// Init initializes the parser; we do this since we need to pull all the latest assets that are tradable.
// Currently we use Binance.
func Init(ctx context.Context) error {
	// Build a list of asset pairs that have futures trading enabled.
	var (
		assetPairs []*binanceproto.AssetPair
		err        error
	)
	for i := 0; i < 3; i++ {
		rsp, retryErr := (&binanceproto.ListAllAssetPairsRequest{}).Send(context.Background()).Response()
		if retryErr != nil {
			slog.Info(ctx, "Failed to fetch all asset pairs on binance. Trying again...", err)
			time.Sleep(500 * time.Millisecond)
			err = retryErr
			continue
		}

		assetPairs = rsp.AssetPairs
		break
	}

	if len(assetPairs) == 0 && err != nil {
		return gerrors.Augment(err, "failed_to_init_parser.failed_to_fetch_asset_pairs_from_binance", nil)
	}

	for _, assetPair := range assetPairs {
		// We may have some inconsistencies here since we're using all symbols; since some symbols
		// May not be actively traded on Binance.
		binanceAssetPairs[strings.ToLower(assetPair.BaseAsset)] = true
	}

	slog.Info(context.Background(), "Fetched all binance asset pairs for satoshi parser; total: %v", len(binanceAssetPairs))
	return nil
}

// Parse ...
func Parse(ctx context.Context, identifier, content string, m *discordgo.MessageCreate) (*tradeengineproto.Trade, error) {
	parser, ok := getParserByIdentifier(identifier)
	if !ok {
		return nil, gerrors.FailedPrecondition("failed_to_parse.parser_does_not_exist", nil)
	}

	cleanedContent := cleanContent(content)

	return parser.Parse(ctx, cleanedContent, m)
}

func cleanContent(content string) string {
	// Remove the dollar sign.
	c := strings.ReplaceAll(content, "$", "")

	// Remove commas.
	c = strings.ReplaceAll(c, ",", "")

	// TODO: Remove Attachments
	// TODO: Remove --- old ---
	// TODO: Remove --- in reply too ---

	// Normalize
	c = strings.ToLower(c)

	return c
}
