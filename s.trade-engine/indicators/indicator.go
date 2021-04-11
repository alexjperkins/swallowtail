package indicators

import (
	"context"
	"fmt"
	"time"

	"swallowtail/libraries/structures/window"

	"github.com/monzo/slog"
)

var (
	BuySignal    = "BUY"
	SellSignal   = "SELL"
	NoSideSignal = "NOSIDE"
)

// Cache is a map of moving windows, used to hold prior info based on dimensions requried for an indicator
type Cache map[string]*window.MovingWindow

// Dimensions is an array of strings, where each string is a dimension; such as volume, low price, high price etc
type Dimensions []string

type Signal struct {
	Side       string
	Reason     string
	Timestamp  time.Time
	Confidence float32
	Metadata   map[string]string
}

type SignalGenerator func(ctx context.Context, event interface{}, output chan *Signal, cache Cache)

type Indicator struct {
	ID string
	// To be used as a hook to push events to
	Input chan interface{}
	// Channel to receive SSL signals
	Signals chan *Signal
	done    chan struct{}
	cache   Cache
}

func NewIndicator(id string, bufSize int, period int, dimensions Dimensions) *Indicator {
	cache := map[string]*window.MovingWindow{}
	for _, dimension := range dimensions {
		cache[dimension] = window.NewMovingWindow(period)
	}

	return &Indicator{
		Input:   make(chan interface{}, bufSize),
		Signals: make(chan *Signal, bufSize),
		done:    make(chan struct{}, 1),
		cache:   cache,
	}
}

func (i *Indicator) Stop() {
	i.done <- struct{}{}
}

func IndicatorFactory(ctx context.Context, id string, bufSize int, period int, signalGenerator SignalGenerator, dimensions Dimensions) *Indicator {
	i := NewIndicator(id, bufSize, period, dimensions)
	assertDimensions(i.cache, dimensions)
	go func() {
		defer func() {
			close(i.Signals)
			slog.Info(ctx, "Closing SSL Channel Indicator.")
		}()
		for {
			select {
			case e := <-i.Input:
				signalGenerator(ctx, e, i.Signals, i.cache)
			case <-i.done:
				return
			case <-ctx.Done():
				i.Stop()
				return
			}
		}
	}()
	return i
}

func assertDimensions(cache Cache, dimensions Dimensions) {
	for _, d := range dimensions {
		if _, ok := cache[d]; !ok {
			panic(fmt.Sprintf("cache missing dimension: %v", d))
		}
	}
}
