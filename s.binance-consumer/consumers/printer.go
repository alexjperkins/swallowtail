package consumers

import (
	"fmt"
	"reflect"
	"strconv"

	"swallowtail/s.binance/domain"

	"github.com/fatih/color"
	"github.com/monzo/slog"
)

var (
	instruments = []string{}
)

func init() {
	red := color.New(color.FgRed).SprintfFunc()
	green := color.New(color.FgGreen).SprintfFunc()
	for i, ticker := range instruments {
		i := i
		register(ticker, func() chan interface{} {
			c := make(chan interface{}, 16)
			go func() {
				// we continue until upstream closes the channel
				var color func(format string, a ...interface{}) string
				for val := range c {
					v := reflect.Indirect(reflect.ValueOf(val)).Interface().(*domain.BinanceKlineEvent)
					switch greaterThan(v.Data.OpenPrice, v.Data.ClosePrice) {
					case true:
						color = red
					case false:
						color = green
					}
					slog.Info(nil, color(
						"[%d]: %d <%s_%s>: o: %v c: %v", i, v.EventTime,
						v.Symbol, v.Data.Interval, v.Data.OpenPrice, v.Data.ClosePrice,
					))
				}
				fmt.Println("Closing channel...")
			}()
			return c
		}(), printerFilter, nil)
	}
}

func greaterThan(opened, closed string) bool {
	so, err := strconv.ParseFloat(opened, 32)
	if err != nil {
		// Best effort
		return true
	}
	sc, err := strconv.ParseFloat(closed, 32)
	if err != nil {
		// Best effort
		return true
	}
	return so > sc
}

func printerFilter(e interface{}) (interface{}, bool) {
	return e, true
}
