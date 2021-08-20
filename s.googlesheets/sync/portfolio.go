package sync

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/monzo/slog"
	"github.com/monzo/terrors"

	"swallowtail/libraries/gerrors"
	"swallowtail/s.googlesheets/dao"
	"swallowtail/s.googlesheets/domain"
	"swallowtail/s.googlesheets/spreadsheet"
	"swallowtail/s.googlesheets/templates"
)

var (
	defaultMaxJitterRange = 30.0
	defaultJitterUnit     = time.Second
)

func init() {
	register("portfolio-syncer", &PortfolioSyncer{
		Spreadsheet: &spreadsheet.Portfolio{},
		change:      make(chan struct{}, 1),
	})
}

// GoogleSheetsPortfolioSyncer ...
type PortfolioSyncer struct {
	// internal spreadsheet specifically for the portfolio syncer.
	Spreadsheet *spreadsheet.Portfolio
	// The entry price increase for which warrants a pager.
	IncreaseAmountsPagers []float64
	Sheets                []*domain.Googlesheet
	sheetsMu              sync.RWMutex
	change                chan struct{}
}

// Sync ...
func (p *PortfolioSyncer) Sync(ctx context.Context) {
	for {
		// This is slightly inefficient since we shutdown and resync all sheets once we get a change recieved
		// but we do avoid having to maintain state in a concurrent environment.
		childCtx, cancel := context.WithCancel(ctx)
		for _, sheet := range p.Sheets {
			sheet := sheet
			go p.sync(childCtx, sheet.UserID, sheet.SpreadsheetID, sheet.SheetID)
		}

		select {
		case <-p.change:
			slog.Debug(ctx, "Portfolio Syncer: change notifcation received")
			cancel()

			// Sleep to allow for a graceful shutdown of goroutines.
			time.Sleep(10 * time.Second)
		case <-ctx.Done():
			cancel()
			return
		}
	}
}

// Refresh ...
func (p *PortfolioSyncer) Refresh(ctx context.Context) error {
	// Initial load from our persistence storage.
	var err error
	for i := 0; i < 5; i++ {
		e := p.refresh(ctx)
		if e == nil {
			break
		}

		slog.Debug(ctx, "%v) Initial refresh attempt error: %v", i, e)

		multierror.Append(err, e)
		time.Sleep(5 * time.Second)
	}

	// We've tried 5 times on start & we cannot load from dao; let's fail.
	if err != nil {
		return terrors.Augment(err, "Failed to perform inital loading of sheets data; with 5 retries", nil)
	}

	go func() {
		// Ideally this should be a consumer of async events; but we don't yet have that infra structure in place.
		// So we have to poll rather than push.
		t := time.NewTicker(1 * time.Minute)
		for {
			select {
			case <-t.C:
				// Best effort
				p.refresh(ctx)
			case <-ctx.Done():
			}
		}
	}()

	return nil
}

func (p *PortfolioSyncer) refresh(ctx context.Context) error {
	ss, err := dao.ListSheetsByType(ctx, templates.PortfolioSheetType)
	switch {
	case terrors.Is(err, "not_found.no-googlesheets-registered-with-this-type"):
		// No sheets
		return nil
	case err != nil:
		return terrors.Augment(err, "Failed to refresh portfolio syncer internal list of sheets", nil)
	}

	if len(ss) == len(p.Sheets) {
		// Since we can't delete/mutate sheets, if the length of the lists are the same, then we know we don't have
		// any changes & we don't have to take a lock.
		return nil
	}

	select {
	case p.change <- struct{}{}:
	default:
		// If we're blocked on sending; then we haven't yet refreshed on the last change.
		// we can skip.
	}

	p.sheetsMu.Lock()
	defer p.sheetsMu.Unlock()
	p.Sheets = ss
	return nil
}

func (p *PortfolioSyncer) sync(ctx context.Context, userID, spreadsheetID, sheetID string) {
	// Add jitter; this prevents us trying to sync everything at once.
	// our exchange client is cached; so we don't need to really worry about rate limiting.
	rand.Seed(time.Now().UnixNano())
	sleepFor := time.Duration(rand.Float64()*defaultMaxJitterRange) * defaultJitterUnit
	time.Sleep(sleepFor)

	t := time.NewTicker(1 * time.Minute)
	for {
		select {
		case <-t.C:
			// Fetch rows
			rows, err := p.Spreadsheet.Rows(ctx, spreadsheetID, sheetID)
			if err != nil {
				if gerr := pageAccount(ctx, userID, fmt.Sprintf("Failed to parse row: Error %v", err.Error()), spreadsheetID); gerr != nil {
					slog.Error(ctx, "Failed to page account: %v", gerr)
				}
				continue
			}

			// Fetch latest prices from our exchange.
			wg := sync.WaitGroup{}
			for i, row := range rows {
				i, row := i, row
				wg.Add(1)
				go func() {
					defer wg.Done()

					assetPair, ok := isValidAssetPairOrConvert(row.AssetPair)
					if !ok {
						// Page user; we don't have a valid asset pair.
						if gerr := pageAccount(ctx, userID, fmt.Sprintf("Invalid asset pair: row %v, asset pair: %s", row.Index, row.AssetPair), spreadsheetID); gerr != nil {
							slog.Error(ctx, "Failed to page account: %v", gerr)
						}
						return
					}

					// Fetch latest price.
					rsp, err := getLatestPriceBySymbol(ctx, row.Ticker, assetPair)
					switch {
					case
						gerrors.Is(err, gerrors.ErrDeadlineExceeded):
						// This is fine; we can retry on the next attempt.
						return
					case err != nil:
						// We failed lets page the user.
						if gerr := pageAccount(ctx, userID, fmt.Sprintf("Failed to retrieve price for asset with ticker: %s row: %v", row.Ticker, row.Index), spreadsheetID); gerr != nil {
							slog.Warn(ctx, "Failed to page account: %v", gerr)
						}
						return
					}

					row.CurrentPrice = float64(rsp.GetLatestPrice())

					// Refresh our row with our latest current price & reset our list.
					row.Refresh()
					rows[i] = row
				}()
			}
			wg.Wait()

			// Upate all rows with our latest price.
			err = p.Spreadsheet.UpdateRows(ctx, spreadsheetID, sheetID, rows)
			if err != nil {
				slog.Info(ctx, "Failed to upload googlesheet row", map[string]string{
					"user_id":  userID,
					"sheet_id": sheetID,
				})
				continue
			}

			// Calculate & update historical PNL
			var historicalPNL float64
			h, err := p.Spreadsheet.TradeHistory(ctx, spreadsheetID, sheetID)
			switch {
			case err != nil:
				slog.Info(ctx, "Failed to read historical data", map[string]string{
					"user_id":   userID,
					"sheet_id":  sheetID,
					"error_msg": err.Error(),
				})
			default:
				historicalPNL = calculateHistoricalPNL(h)
				if err := p.Spreadsheet.UpdateTradeHistory(ctx, spreadsheetID, sheetID, h); err != nil {
					// Best effort
					slog.Info(ctx, "Failed to re-upload refreshed historical trades: %v", err)
				}
			}

			// Update metadata
			m, err := p.Spreadsheet.Metadata(ctx, spreadsheetID, sheetID)
			if err != nil {
				if gerr := pageAccount(ctx, userID, "Apologies, I wasn't able to parse the metadata in your portfolio tracker; please check it's correct", spreadsheetID); gerr != nil {
					slog.Error(ctx, "Failed to parse account: %v", gerr)
				}
				// We can't go any further; lets skip.
				continue
			}

			// Calculate our total worth
			m.TotalPNL = calculateTotalPNL(rows) + historicalPNL

			m.TotalWorth, err = calculateTotalWorth(ctx, rows, m.AssetPair)
			if err != nil {
				slog.Error(ctx, "Failed to update googlesheet metadata", map[string]string{
					"error_msg": err.Error(),
				})
			}

			err = p.Spreadsheet.UpdateMetadata(ctx, spreadsheetID, sheetID, m)
			if err != nil {
				slog.Info(ctx, "Failed to update googlesheet metadata", map[string]string{
					"user_id":  userID,
					"sheet_id": sheetID,
				})
				continue
			}
			slog.Info(ctx, "Updated googlesheet metadata", map[string]string{
				"user_id":  userID,
				"sheet_id": sheetID,
			})

		case <-ctx.Done():
			return
		}
	}
}
