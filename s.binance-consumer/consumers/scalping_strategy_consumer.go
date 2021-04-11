package consumers

import (
	"context"
	"fmt"
	"reflect"
	"swallowtail/s.binance-consumer/domain"
	"swallowtail/s.binance-consumer/strategies"

	"github.com/monzo/slog"
)

func init() {
	tickers := []string{"ethusdt", "btcusdt"}
	for _, ticker := range tickers {
		// Register SSL chan
		sslInstrument := buildSSLInstrument(ticker)
		sslCh := make(chan interface{}, 16)
		register(sslInstrument, sslCh, sslFilter, nil)

		// Regsiter EMA chan
		// TODO

		// Register dwarf chan
		// TODO

		// Start strategy on consumer
		go strategies.NewScalpingStrategy(context.Background(), ticker, sslCh).Run()
	}

}

func buildSSLInstrument(ticker string) string {
	return fmt.Sprintf("%s@kline_1m", ticker)
}

func sslFilter(e interface{}) (interface{}, bool) {
	// We expect the event to be a binance kline event, but we only care for events
	// that are closed;
	ce := reflect.Indirect(reflect.ValueOf(e)).Interface().(*domain.BinanceKlineEvent)
	slog.Debug(nil, "ticker -> %v, closed -> %v, price ->%v", ce.Data.Symbol, ce.Data.IsKlineClosed, ce.Data.ClosePrice)
	return ce, ce.Data.IsKlineClosed
}
