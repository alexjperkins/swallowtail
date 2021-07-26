package domain

import (
	"context"
	"fmt"
	"reflect"
	"swallowtail/libraries/util"
	"time"

	"github.com/monzo/slog"
)

// Googlesheet ...
type Googlesheet struct {
	GooglesheetID     string    `pb:"googlesheet_id"`
	SpreadsheetID     string    `pb:"spreadsheet_id"`
	SheetID           string    `pb:"sheet_id"`
	UserID            string    `pb:"user_id"`
	WithPagerOnError  bool      `pb:"with_pager_on_error"`
	WithPagerOnTarget bool      `pb:"with_pager_on_target"`
	Created           time.Time `pb:"created"`
	Updated           time.Time `pb:"updated"`
	Active            bool      `pb:"active"`
}

// TODO: move the below

type PortfolioMetadata struct {
	TotalPNL   float64
	TotalWorth float64
	AssetPair  string
}

type PortfolioRow struct {
	Index         int
	Ticker        string
	AssetPair     string
	AverageEntry  float64
	Bought        float64
	Amount        float64
	CurrentPrice  float64
	CurrentValue  float64
	PNL           float64
	PNLPercentage float64
	Target        float64
}

type HistoricalTradeRow struct {
	Index          int
	Ticker         string
	AssetPair      string
	BoughtFor      float64
	Amount         float64
	Sold           float64
	SoldPercentage float64
	PNL            float64
}

func (pr *PortfolioRow) ToArray() []interface{} {
	r := []interface{}{}
	v := reflect.ValueOf(pr)
	if v.IsZero() {
		return []interface{}{}
	}
	e := v.Elem()
	n := e.NumField()
	for i := 1; i < n; i++ {
		f := e.Field(i)
		switch f.Kind() {
		case reflect.String:
			r = append(r, f.String())
		case reflect.Int:
			s := fmt.Sprintf("%v", f.Int())
			r = append(r, s)
		case
			reflect.Float64,
			reflect.Float32:
			s, err := util.FormatPriceAsString(f.Float())
			if err != nil {
				slog.Error(context.TODO(), "Failed to format current price as string")
				s = ""
			}
			r = append(r, s)
		default:
			panic("Unhandled field type Portfolio Row")
		}
	}
	return r
}

func (pr *PortfolioRow) Refresh() {
	pr.CurrentValue = pr.CurrentPrice * pr.Amount
	pr.PNL = pr.CurrentValue - pr.Bought
	pr.PNLPercentage = calcPNLPerc(pr.CurrentPrice, pr.AverageEntry)
}

func (pr *PortfolioRow) WithTarget() bool {
	// We define no target, if it's set to the default value; this makes sense
	// since what rationale person sets an investment target of zero?
	return pr.Target == 0.0
}

func (ht *HistoricalTradeRow) Refresh() {
	ht.PNL = ht.Amount * (ht.Sold - ht.BoughtFor) * ht.SoldPercentage * 0.01
}

func (ht *HistoricalTradeRow) ToArray() []interface{} {
	r := []interface{}{}
	v := reflect.ValueOf(ht)
	if v.IsZero() {
		return []interface{}{}
	}
	e := v.Elem()
	n := e.NumField()
	for i := 1; i < n; i++ {
		f := e.Field(i)
		switch f.Kind() {
		case reflect.String:
			r = append(r, f.String())
		case reflect.Int:
			s := fmt.Sprintf("%v", f.Int())
			r = append(r, s)
		case
			reflect.Float64,
			reflect.Float32:
			s, err := util.FormatPriceAsString(f.Float())
			if err != nil {
				slog.Error(context.TODO(), "Failed to format current price as string")
				s = ""
			}
			r = append(r, s)
		default:
			panic("Unhandled field type Portfolio Row")
		}
	}
	return r
}

func calcPNLPerc(currentPrice, boughtAt float64) float64 {
	if boughtAt == 0.0 {
		return 0.0
	}
	v := ((currentPrice / boughtAt) - 1) * 100
	if v < 0.01 && v > 0.0 {
		return 0.0
	}
	return v
}
