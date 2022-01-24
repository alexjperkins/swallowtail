package exchangeinfo

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/monzo/slog"
	"github.com/tidwall/gjson"

	"swallowtail/libraries/gerrors"
	"swallowtail/s.binance/client"
)

type FilterType int

const (
	FilterTypePrice FilterType = iota + 1
	FilterTypeLotSize
	FilterTypeMarketLotSize
)

func (f FilterType) String() string {
	switch f {
	case FilterTypePrice:
		return "PRICE_FILTER"
	case FilterTypeLotSize:
		return "LOT_SIZE"
	case FilterTypeMarketLotSize:
		return "MARKET_LOT_SIZE"
	default:
		return ""
	}
}

type SymbolData struct {
	ContractType string
	Symbol       string
	Pair         string
	BaseAsset    string

	MinPrice string
	MaxPrice string
	TickSize string

	LotSize           string
	MinQuantity       string
	MarketLotSize     string
	MarketMinQuantity string
}

var (
	symbolInformation = map[string]*SymbolData{}
	mu                sync.RWMutex
)

// Init initializes the exchange information required for this service.
func Init(ctx context.Context) error {
	if err := gatherExchangeInfo(ctx); err != nil {
		return err
	}

	slog.Info(ctx, "Gathered required futures exchange information: #%d assets", len(symbolInformation))

	// Start our refresh loop.
	go refresh(ctx)

	return nil
}

func refresh(ctx context.Context) {
	t := time.NewTicker(23 * time.Hour)
	for {
		select {
		case <-t.C:
			if err := gatherExchangeInfo(ctx); err != nil {
				slog.Error(ctx, "Failed to refresh binance exchange info: Error: %v", err)
				continue
			}
			slog.Info(ctx, "Refreshed binance exchange info")
		case <-ctx.Done():
			return
		}
	}
}

func gatherExchangeInfo(ctx context.Context) error {
	var (
		rsp *client.GetFuturesExchangeInfoResponse
		err error
	)
	for i := 0; i < 3; i++ {
		r, e := client.GetFuturesExchangeInfo(ctx, &client.GetFuturesExchangeInfoRequest{})
		if e != nil {
			multierror.Append(err, e)
			slog.Trace(ctx, "Failed to gather exchangeinfo, attempt [%v]; retrying...", i)
		}

		rsp = r
		break
	}

	if err != nil {
		return gerrors.Augment(err, "failed_to_init_exchange_info", nil)
	}
	if rsp == nil {
		return gerrors.Augment(err, "failed_to_init_exchange_info.empty_response", nil)
	}

	mu.Lock()
	defer mu.Unlock()

	for _, s := range rsp.Symbols {
		var (
			minPrice      string
			maxPrice      string
			tickSize      string
			lotSize       string
			minQty        string
			marketLotSize string
			marketMinQty  string
		)

		b, err := json.Marshal(s.Filters)
		if err != nil {
			return gerrors.Augment(err, "failed_to_init_exchange_info.bad_filter", map[string]string{
				"symbol": s.Symbol,
			})
		}

		ff := gjson.ParseBytes(b)
		for _, f := range ff.Array() {
			if ft := f.Get("filterType"); ft.Exists() {
				switch {
				case FilterTypePrice.String() == ft.String():
					switch v := f.Get("maxPrice"); {
					case !v.Exists() || v.String() == "":
						slog.Error(ctx, "Failed to parse ; max price: %s: %v", s.Symbol, err)
					default:
						maxPrice = v.String()
					}

					switch v := f.Get("minPrice"); {
					case !v.Exists() || v.String() == "":
						slog.Error(ctx, "Failed to parse ; min price: %s: %v", s.Symbol, err)
					default:
						minPrice = v.String()
					}

					switch v := f.Get("tickSize"); {
					case !v.Exists() || v.String() == "":
						slog.Error(ctx, "Failed to parse ; tick size: %s: %v", s.Symbol, err)
					default:
						tickSize = v.String()
					}
				case FilterTypeLotSize.String() == ft.String():
					switch v := f.Get("stepSize"); {
					case !v.Exists() || v.String() == "":
						slog.Error(ctx, "Failed to parse ; lot size: %s: %v", s.Symbol, err)
					default:
						lotSize = v.String()
					}

					switch v := f.Get("minQty"); {
					case !v.Exists() || v.String() == "":
						slog.Error(ctx, "Failed to parse ; min qty: %s: %v", s.Symbol, err)
					default:
						minQty = v.String()
					}
				case FilterTypeMarketLotSize.String() == ft.String():
					switch v := f.Get("stepSize"); {
					case !v.Exists() || v.String() == "":
						slog.Error(ctx, "Failed to parse ; market lot size: %s: %v", s.Symbol, err)
					default:
						marketLotSize = v.String()
					}

					switch v := f.Get("minQty"); {
					case !v.Exists() || v.String() == "":
						slog.Error(ctx, "Failed to parse ; market min qty: %s: %v", s.Symbol, err)
					default:
						marketMinQty = v.String()
					}
				}
			}

			// We can continue if the filter type doesn't exist in the object.
			continue
		}

		d := &SymbolData{
			BaseAsset:         s.BaseAsset,
			ContractType:      s.ContractType,
			Pair:              s.Pair,
			Symbol:            s.Symbol,
			MinPrice:          minPrice,
			MaxPrice:          maxPrice,
			TickSize:          tickSize,
			LotSize:           lotSize,
			MarketLotSize:     marketLotSize,
			MinQuantity:       minQty,
			MarketMinQuantity: marketMinQty,
		}

		symbolInformation[strings.ToLower(s.Symbol)] = d
	}

	return nil
}

// GetBaseAssetQuantityPrecision returns the base asset quantity precision given the base asset.
func GetBaseAssetQuantityPrecision(baseAsset string, isMarketOrder bool) (float64, bool, error) {
	mu.RLock()
	defer mu.RUnlock()

	v, ok := symbolInformation[strings.ToLower(baseAsset)]
	if !ok {
		return 0, false, nil
	}

	var vq string
	switch {
	case isMarketOrder:
		vq = v.MarketLotSize
	default:
		vq = v.LotSize
	}

	if vq == "" {
		return 0, false, nil
	}

	f, err := strconv.ParseFloat(vq, 64)
	if err != nil {
		return 0.0, false, gerrors.Augment(err, "failed_to_parse_lot_size", map[string]string{
			"lot_size": vq,
		})
	}

	return f, true, nil
}

// GetBaseAssetPricePrecision returns the base asset price precision given the base asset.
func GetBaseAssetPricePrecision(baseAsset string) (float64, bool, error) {
	mu.RLock()
	defer mu.RUnlock()

	v, ok := symbolInformation[strings.ToLower(baseAsset)]
	if !ok {
		return 0, false, nil
	}

	if v.TickSize == "" {
		return 0, false, nil
	}

	f, err := strconv.ParseFloat(v.TickSize, 64)
	if err != nil {
		return 0.0, false, gerrors.Augment(err, "failed_to_parse_tick_size", map[string]string{
			"tick_size": v.TickSize,
		})
	}

	return f, true, nil
}

// GetBaseAssetMinQty returns the base asset minimum quantity.
func GetBaseAssetMinQty(baseAsset string, isMarketOrder bool) (float64, bool, error) {
	mu.RLock()
	defer mu.RUnlock()

	v, ok := symbolInformation[strings.ToLower(baseAsset)]
	if !ok {
		return 0, false, nil
	}

	var vq string
	switch {
	case isMarketOrder:
		vq = v.MarketMinQuantity
	default:
		vq = v.MinQuantity
	}

	vf, err := strconv.ParseFloat(vq, 64)
	if err != nil {
		return 0, false, gerrors.Augment(err, "failed_to_get_base_asset_min_qty.bad_value", map[string]string{
			"value": vq,
		})
	}

	return vf, true, nil
}
