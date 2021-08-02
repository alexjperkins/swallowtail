package sync

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"google.golang.org/grpc"

	"github.com/hashicorp/go-multierror"
	"github.com/monzo/slog"
	"github.com/monzo/terrors"

	accountproto "swallowtail/s.account/proto"
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
}

// Sync ...
func (p *PortfolioSyncer) Sync(ctx context.Context) error {
	sheets, err := dao.ListSheetsByType(ctx, templates.PortfolioSheetType)
	if err != nil {
		return terrors.Augment(err, "Failed to read sheets by type", map[string]string{
			"sheet_type": templates.PortfolioSheetType.String(),
		})
	}
	for _, sheet := range sheets {
		sheet := sheet
		go p.sync(ctx, sheet.UserID, sheet.SheetID)
	}
	return nil
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
		multierror.Append(err, e)
		time.Sleep(30 * time.Second)
	}

	// We've tried 5 times on start & we cannot load from dao; let's fail.
	if err != nil {
		return terrors.Augment(err, "Failed to perform inital loading of sheets data; with 5 retries", nil)
	}

	go func() {
		// Ideally this should be a consumer of async events; but we don't yet have that infra structure in place.
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
	if err != nil {
		return terrors.Augment(err, "Failed to refresh portfolio syncer internal list of sheets", nil)
	}
	if len(ss) == len(p.Sheets) {
		// Since we can't mutate, if the lenght of the lists are the same, then we know we don't have
		// any changes & we don't have to take a lock.
		return nil
	}

	p.sheetsMu.Lock()
	defer p.sheetsMu.Unlock()
	p.Sheets = ss
	return nil
}

func (p *PortfolioSyncer) sync(ctx context.Context, userID, sheetID string) {
	// Add jitter; this prevents us trying to sync everything at once.
	// our exchange client is cached; so we don't need to really worry about rate limiting.
	rand.Seed(time.Now().UnixNano())
	sleepFor := time.Duration(rand.Float64()*defaultMaxJitterRange) * defaultJitterUnit
	time.Sleep(sleepFor)

	t := time.NewTicker(1 * time.Minute)
	for {
		select {
		case <-t.C:
			// s.account client
			ac, aCloser, err := accountClient(ctx)
			if err != nil {
				slog.Error(ctx, "Failed to connect to s.account: %v", err)
				continue
			}
			defer aCloser()

			// Fetch rows
			rows, err := p.Spreadsheet.Rows(ctx, sheetID)
			if err != nil {
				if _, err := (ac.PageAccount(ctx, &accountproto.PageAccountRequest{
					Content:  fmt.Sprintf("Failed to parse rows, please check: %v", err.Error()),
					UserId:   userID,
					Priority: accountproto.PagerPriority_LOW,
				})); err != nil {
					slog.Error(ctx, "Failed to page account: %v", err)
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
						if _, err := (ac.PageAccount(ctx, &accountproto.PageAccountRequest{
							UserId:   userID,
							Content:  fmt.Sprintf("invalid asset pair: %s", row.AssetPair),
							Priority: accountproto.PagerPriority_LOW,
						})); err != nil {
							slog.Error(ctx, "Failed to page account: %v", err)
						}
						return
					}

					// Fetch latest price.
					rsp, err := getLatestPriceBySymbol(ctx, row.Ticker, assetPair)
					if err != nil {
						slog.Error(ctx, "Failed to get latest price by symbol: %v", err)
					}

					latestPrice := rsp.LatestPrice

					switch {
					case
						terrors.Is(err, terrors.ErrRateLimited),
						terrors.Is(err, terrors.ErrTimeout):
						// This is fine; we can retry on the next attempt.
						return

					case err != nil:
						// We failed lets page the user.
						if _, err := (ac.PageAccount(ctx, &accountproto.PageAccountRequest{
							UserId:  userID,
							Content: fmt.Sprintf("failed to retrieve price for asset with ticker: %s, please check that it is correct.", row.Ticker),
						})); err != nil {
							slog.Warn(ctx, "Failed to page account: %v", err)
						}
						return
					}

					row.CurrentPrice = float64(latestPrice)

					// Refresh our row with our latest current price & reset our list.
					row.Refresh()
					rows[i] = row
				}()
			}
			wg.Wait()

			// Upate all rows with our latest price.
			err = p.Spreadsheet.UpdateRows(ctx, sheetID, rows)
			if err != nil {
				slog.Info(ctx, "Failed to upload googlesheet row", map[string]string{
					"user_id":  userID,
					"sheet_id": sheetID,
				})
				continue
			}

			// Calculate & update historical PNL
			var historicalPNL float64
			h, err := p.Spreadsheet.TradeHistory(ctx, sheetID)
			switch {
			case err != nil:
				slog.Info(ctx, "Failed to read historical data", map[string]string{
					"user_id":   userID,
					"sheet_id":  sheetID,
					"error_msg": err.Error(),
				})
			default:
				historicalPNL = calculateHistoricalPNL(h)
				if err := p.Spreadsheet.UpdateTradeHistory(ctx, sheetID, h); err != nil {
					// Best effort
					slog.Info(ctx, "Failed to re-upload refreshed historical trades.")
				}
			}

			// Update metadata
			m, err := p.Spreadsheet.Metadata(ctx, sheetID)
			if err != nil {
				if _, err := (ac.PageAccount(ctx, &accountproto.PageAccountRequest{
					Content: "Apologies, I wasn't able to parse the metadata in your portfolio tracker; please check it's correct",
					UserId:  userID,
				})); err != nil {
					slog.Error(ctx, "Failed to parse account: %v", err)
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

			err = p.Spreadsheet.UpdateMetadata(ctx, sheetID, m)
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

func accountClient(ctx context.Context) (client accountproto.AccountClient, closer func() error, err error) {
	conn, err := grpc.DialContext(ctx, "swallowtail-s-account:8000", grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}
	return accountproto.NewAccountClient(conn), conn.Close, nil
}
