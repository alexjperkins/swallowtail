package sync

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"swallowtail/libraries/ttlcache"
	"swallowtail/s.googlesheets/domain"
	"swallowtail/s.googlesheets/spreadsheet"
	"sync"
	"time"

	"github.com/monzo/slog"
	"github.com/monzo/terrors"
)

var (
	defaultMaxJitterRange = 30.0
	defaultJitterUnit     = time.Second
	defaultPagerGif       = "https://tenor.com/view/cynical-pepe-cynical-p-pepe-the-frog-frog-cynical-gif-14037546"

	validAssetPairs = map[string]bool{
		"USDT": true,
		"USD":  true,
		"BTC":  true,
		"ETH":  true,
		"GBP":  true,
	}
	validAssetPairMtx sync.RWMutex
)

type ExchangeClient interface {
	GetPrice(ctx context.Context, symbol, assetPair string) (float64, error)
	Ping(ctx context.Context) bool
}

func NewGoogleSheetsPorfolioSyncer(
	googleSpreadsheet *spreadsheet.GoogleSheetPortfolio,
	exchangeClient ExchangeClient,
	interval time.Duration,
	done chan struct{},
	withJitter bool,
) *GoogleSheetsPortfolioSyncer {
	return &GoogleSheetsPortfolioSyncer{
		Spreadsheet: googleSpreadsheet,
		ec:          exchangeClient,
		interval:    interval,
		done:        done,
		withJitter:  withJitter,
	}
}

// GoogleSheetsPortfolioSyncer
type GoogleSheetsPortfolioSyncer struct {
	Spreadsheet *spreadsheet.GoogleSheetPortfolio
	ec          ExchangeClient
	interval    time.Duration
	done        chan struct{}
	withJitter  bool
}

func (p *GoogleSheetsPortfolioSyncer) Start(ctx context.Context) {
	for i, sheetID := range p.Spreadsheet.SheetIDs {
		slog.Info(ctx, "Starting portfolio sync: [%v] %s %s", i, p.Spreadsheet.Owner.Name, sheetID)
		go p.sync(ctx, sheetID)
	}
}

func (p *GoogleSheetsPortfolioSyncer) sync(ctx context.Context, sheetID string) {
	// Basic cache that stores TTL to stop owners being pinged too often.
	ttlcache := ttlcache.New(time.Duration(12 * time.Hour))

	t := time.NewTicker(p.interval)
	defer slog.Info(ctx, "Closing down google sheets price syncer", nil)

	// Add jitter
	rand.Seed(time.Now().UnixNano())
	if p.withJitter {
		sleepFor := time.Duration(rand.Float64()*defaultMaxJitterRange) * defaultJitterUnit
		slog.Info(ctx, "Adding jitter; sleeping for %v", sleepFor)
		time.Sleep(sleepFor)
	}

	for {
		select {
		case <-t.C:
			wg := sync.WaitGroup{}
			rows, err := p.Spreadsheet.Rows(ctx)
			if err != nil {
				slog.Error(ctx, "Failed to retrieve values", map[string]string{
					"spreadsheet_id": p.Spreadsheet.ID(),
					"error":          err.Error(),
				})
				invalidRowsMsg := fmt.Sprintf("Failed to parse rows, please check: %v", err.Error())
				err := p.Spreadsheet.Owner.Page(ctx, formatPagerMsg(p.Spreadsheet.Owner.Name, "", invalidRowsMsg, 0))
				if err != nil {
					slog.Info(ctx, "Failed to page user", map[string]string{
						"user_id":       p.Spreadsheet.Owner.Name,
						"error_message": err.Error(),
					})
				}
				continue
			}
			// Do I need this?
			mutatedRows := make([]*domain.PortfolioRow, len(rows))
			for i, row := range rows {
				i, row := i, row
				wg.Add(1)
				go func() {
					defer wg.Done()

					validAssetPair, ok := isValidAssetPairOrConvert(row.AssetPair)
					if !ok {
						invalidAssetMsg := fmt.Sprintf("Invalid asset pair: %s", row.AssetPair)
						err := p.Spreadsheet.Owner.Page(ctx, formatPagerMsg(p.Spreadsheet.Owner.Name, row.Ticker, invalidAssetMsg, row.Index))
						if err != nil {
							slog.Info(ctx, "Failed to page user", map[string]string{
								"user_id":       p.Spreadsheet.Owner.Name,
								"error_message": err.Error(),
							})
						}
						return
					}

					row.CurrentPrice, err = p.ec.GetPrice(ctx, row.Ticker, validAssetPair)
					switch {
					case
						terrors.Is(err, terrors.ErrInternalService),
						terrors.Is(err, terrors.ErrRateLimited),
						terrors.Is(err, terrors.ErrTimeout):
						slog.Warn(ctx, "Failed to get price for: %s with error: %v", row.Ticker, err.Error())
						return
					case err != nil:
						slog.Warn(ctx, "Failed to get current price", map[string]string{
							"ticker": row.Ticker,
						})
						// Only page if doesn't exist in the cache
						if _, expired := ttlcache.GetAndRefreshExpiry(p.Spreadsheet.Owner.DiscordID); !expired {

							err := p.Spreadsheet.Owner.Page(ctx, formatPagerMsg(p.Spreadsheet.Owner.Name, row.Ticker, err.Error(), row.Index))
							if err != nil {
								slog.Info(ctx, "Failed to page user", map[string]string{
									"user_id":       p.Spreadsheet.Owner.Name,
									"error_message": err.Error(),
								})
							}
						}
					}
					slog.Info(ctx, fmt.Sprintf("Current Price [%s]: %v", row.Ticker, row.CurrentPrice))
					// Update all row values now that price has been updated.
					row.Refresh()
					mutatedRows[i] = row
				}()
			}
			wg.Wait()

			err = p.Spreadsheet.UpdateRows(ctx, sheetID, mutatedRows)
			if err != nil {
				slog.Info(ctx, "Failed to upload googlesheet row", map[string]string{
					"owner":          p.Spreadsheet.Owner.Name,
					"spreadsheet_id": p.Spreadsheet.ID(),
					"sheet_id":       sheetID,
				})
				continue
			}

			var historicalPNL float64
			h, err := p.Spreadsheet.TradeHistory(ctx)
			if err != nil {
				slog.Info(ctx, "Failed to read historical data", map[string]string{
					"owner":          p.Spreadsheet.Owner.Name,
					"spreadsheet_id": p.Spreadsheet.ID(),
					"sheet_id":       sheetID,
					"error_msg":      err.Error(),
				})
			} else {
				historicalPNL = calculateHistoricalPNL(h)
				if err := p.Spreadsheet.UpdateTradeHistory(ctx, sheetID, h); err != nil {
					// Best effort
					slog.Info(ctx, "Failed to re-upload refreshed historical trades.")
				}
			}

			m, err := p.Spreadsheet.Metadata(ctx)
			if err != nil {
				slog.Info(ctx, "Failed to read googlesheet metadata", map[string]string{
					"owner":          p.Spreadsheet.Owner.Name,
					"spreadsheet_id": p.Spreadsheet.ID(),
					"sheet_id":       sheetID,
				})
				continue
			}

			m.TotalPNL = calculateTotalPNL(mutatedRows) + historicalPNL
			m.TotalWorth = calculateTotalWorth(mutatedRows)

			err = p.Spreadsheet.UpdateMetadata(ctx, sheetID, m)
			if err != nil {
				slog.Info(ctx, "Failed to update googlesheet metadata", map[string]string{
					"owner":          p.Spreadsheet.Owner.Name,
					"spreadsheet_id": p.Spreadsheet.ID(),
					"sheet_id":       sheetID,
				})
				continue
			}
			slog.Info(ctx, "Updated googlesheet metadata", map[string]string{
				"owner":          p.Spreadsheet.Owner.Name,
				"spreadsheet_id": p.Spreadsheet.ID(),
				"sheet_id":       sheetID,
			})
		case <-ctx.Done():
			return
		case <-p.done:
			return
		}
	}
}

func formatPagerMsg(name, ticker, errorMsg string, rowIndex int) string {
	return fmt.Sprintf(
		":wave: Hello there, %s\n```Failed to get price for row: `%v` with ticker `%s`\nError: %v\n```Please check it is correct.\n%s\n",
		name, strconv.Itoa(rowIndex), ticker, errorMsg, defaultPagerGif,
	)
}

func calculateTotalPNL(rows []*domain.PortfolioRow) float64 {
	var total float64
	for _, row := range rows {
		total += row.PNL
	}
	return total
}

func calculateTotalWorth(rows []*domain.PortfolioRow) float64 {
	var total float64
	for _, row := range rows {
		total += row.CurrentValue
	}
	return total
}

// calculateHistoricalPNL iterates through all rows, refreshes them to recalculate the PNL
// and returns the total.
func calculateHistoricalPNL(rows []*domain.HistoricalTradeRow) float64 {
	var total float64
	for _, row := range rows {
		row.Refresh()
		total += row.PNL
	}
	return total
}

func isValidAssetPairOrConvert(assetPair string) (string, bool) {
	validAssetPairMtx.RLock()
	defer validAssetPairMtx.RUnlock()
	_, ok := validAssetPairs[assetPair]
	if !ok {
		return "", false
	}
	if strings.ToLower(assetPair) == "usdt" {
		return "usd", true
	}
	return assetPair, true
}
