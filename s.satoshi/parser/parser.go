package parser

import (
	"context"
	"strings"
	"swallowtail/libraries/gerrors"
	binanceproto "swallowtail/s.binance/proto"
	tradeengineproto "swallowtail/s.trade-engine/proto"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/hashicorp/go-multierror"
	"github.com/monzo/slog"
)

var (
	binanceAssetPairs = map[string]bool{}
)

// TradeParser ...
type TradeParser interface {
	Parse(ctx context.Context, content string, m *discordgo.MessageCreate, actorType tradeengineproto.ACTOR_TYPE) (*tradeengineproto.Trade, error)
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
			time.Sleep(2 * time.Second)
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
func Parse(ctx context.Context, identifier, content string, m *discordgo.MessageCreate, actorType tradeengineproto.ACTOR_TYPE) (*tradeengineproto.Trade, error) {
	parsers, ok := getParsersByIdentifier(identifier)
	if !ok {
		return nil, gerrors.FailedPrecondition("failed_to_parse.parser_does_not_exist", nil)
	}

	cleanedContent := cleanContent(content)

	var mErr error
	for _, parser := range parsers {
		trade, err := parser.Parse(ctx, cleanedContent, m, actorType)
		if err != nil {
			slog.Error(ctx, "Failed to parse trade: %v", err)
			multierror.Append(mErr, err)
			continue
		}

		return trade, nil
	}

	return nil, gerrors.Augment(mErr, "Failed to parse trades using any parser", nil)
}

func cleanContent(content string) string {
	// Remove the dollar sign.
	c := strings.ReplaceAll(content, "$", " ")

	// Remove commas.
	c = strings.ReplaceAll(c, ",", "")

	// Remove tabs & newlines; replace with spaces & then trim & remove excess spaces.
	c = strings.ReplaceAll(c, "\n", " ")
	c = strings.ReplaceAll(c, "\t", " ")
	c = strings.TrimSpace(c)
	c = strings.ReplaceAll(c, "  ", " ")

	// TODO: Remove Attachments
	// TODO: Remove --- in reply too ---

	// Normalize
	c = strings.ToLower(c)

	// Remove any replies.
	if strings.Contains(c, "--- message was a reply to ---") {
		splits := strings.Split(c, "--- message was a reply to ---")
		c = splits[0]
	}

	// Remove all old messages.
	if !strings.Contains(c, "---new---") {
		return c
	}

	// Take everything after "---new---".
	updates := strings.SplitAfter(c, "---new---")
	return updates[len(updates)-1]
}
