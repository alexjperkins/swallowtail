package parser

import (
	"context"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/hashicorp/go-multierror"
	"github.com/monzo/slog"

	"swallowtail/libraries/gerrors"
	"swallowtail/libraries/util"
	binanceproto "swallowtail/s.binance/proto"
	ftxproto "swallowtail/s.ftx/proto"
	tradeengineproto "swallowtail/s.trade-engine/proto"
)

var (
	binanceInstruments = map[string]bool{}
	ftxInstruments     = map[string]bool{}
)

// TradeParser ...
type TradeParser interface {
	Parse(ctx context.Context, content string, m *discordgo.MessageCreate, actorType tradeengineproto.ACTOR_TYPE) (*tradeengineproto.TradeStrategy, error)
}

// Init initializes the parser; we do this since we need to pull all the latest assets that are tradable.
// Currently we use Binance.
func Init(ctx context.Context) error {
	// Fetch Binance info.
	rsp, err := util.Retry(ctx, 5, func(ctx context.Context) (interface{}, error) {
		return (&binanceproto.ListAllAssetPairsRequest{}).Send(context.Background()).Response()
	})
	if err != nil {
		return gerrors.Augment(err, "failed_to_init_parser.fetch_binance_info", nil)
	}

	binanceInfo, ok := rsp.(*binanceproto.ListAllAssetPairsResponse)
	if !ok {
		return gerrors.Augment(err, "failed_to_init_parser.bad_list_all_asset_pairs_binance_response.type", nil)
	}

	if len(binanceInfo.AssetPairs) == 0 {
		return gerrors.Augment(err, "failed_to_init_parser.no_binance_asset_pairs", nil)
	}

	slog.Trace(context.Background(), "Fetched all binance asset pairs for satoshi parser; total: %v", len(binanceInstruments))

	// Fetch FTX info.
	rsp, err = util.Retry(ctx, 5, func(ctx context.Context) (interface{}, error) {
		return (&ftxproto.ListFTXInstrumentsRequest{}).Send(context.Background()).Response()
	})
	if err != nil {
		return gerrors.Augment(err, "failed_to_init_parser.fetch_ftx_info", nil)
	}

	ftxInfo, ok := rsp.(*ftxproto.ListFTXInstrumentsResponse)
	if !ok {
		return gerrors.Augment(err, "failed_to_init_parser.bad_list_instruments_ftx_response.type", nil)
	}

	if len(ftxInfo.Instruments) == 0 {
		return gerrors.Augment(err, "failed_to_init_parser.no_ftx_instruments", nil)
	}

	slog.Trace(context.Background(), "Fetched all ftx instruments for satoshi parser; total: %v", len(ftxInstruments))

	// Add to internal instruments cache.
	for _, assetPair := range binanceInfo.AssetPairs {
		// We may have some inconsistencies here since we're using all symbols; since some symbols
		// May not be actively traded on Binance.
		binanceInstruments[strings.ToLower(assetPair.BaseAsset)] = true
	}

	for _, instrument := range ftxInfo.Instruments {
		ftxInstruments[strings.ToLower(instrument.Symbol)] = true
	}

	return nil
}

// Parse ...
func Parse(ctx context.Context, identifier, content string, m *discordgo.MessageCreate, actorType tradeengineproto.ACTOR_TYPE) (*tradeengineproto.TradeStrategy, error) {
	parsers, ok := getParsersByIdentifier(identifier)
	if !ok {
		return nil, gerrors.FailedPrecondition("failed_to_parse.parser_does_not_exist", nil)
	}

	cleanedContent := cleanContent(content)

	var mErr error
	for _, parser := range parsers {
		trade, err := parser.Parse(ctx, cleanedContent, m, actorType)
		if err != nil {
			slog.Error(ctx, "Failed to parse trade: %v %v", err, cleanedContent)
			mErr = multierror.Append(mErr, err)
			continue
		}

		return trade, nil
	}

	if mErr != nil {
		return nil, gerrors.Augment(mErr, "Failed to parse trades using any parser", nil)
	}

	slog.Error(ctx, "Invalid parse of message; empty error & empty trade.")
	return nil, gerrors.FailedPrecondition("invalid_state.no_trade_parsed_without_error", nil)
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
	// Remove all odd & even number of spaces.
	c = strings.ReplaceAll(c, "   ", " ")
	c = strings.ReplaceAll(c, "  ", " ")

	// Remove discord tags.
	c = strings.ReplaceAll(c, "@", "")
	c = strings.ReplaceAll(c, "​", "")

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
