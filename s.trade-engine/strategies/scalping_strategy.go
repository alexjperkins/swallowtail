package strategies

import (
	"context"
	"swallowtail/s.trade-engine/indicators"
	"sync"

	"github.com/monzo/slog"
)

var (
	ScalpingStrategyID = "scaling_strategy"
)

type StrategyInput chan interface{}

type ScalpingStrategy struct {
	ID            string
	Ticker        string
	inputs        map[string]StrategyInput
	inputsMtx     sync.RWMutex
	indicators    map[string]*indicators.Indicator
	indicatorsMtx sync.RWMutex
	cancelFunc    context.CancelFunc
	done          chan struct{}
}

func NewScalpingStrategy(ctx context.Context, ticker string, sslInput StrategyInput) *ScalpingStrategy {
	ctx, cancelFunc := context.WithCancel(ctx)
	return &ScalpingStrategy{
		ID:     ScalpingStrategyID,
		Ticker: ticker,
		inputs: map[string]StrategyInput{
			indicators.SSLInput: sslInput,
		},
		indicators: map[string]*indicators.Indicator{
			indicators.SSLIndicatorID: indicators.SSLChannelIndicator(ctx, 16, 10),
		},
		indicatorsMtx: sync.RWMutex{},
		cancelFunc:    cancelFunc,
		done:          make(chan struct{}, 1),
	}
}

func (ss *ScalpingStrategy) Run() {
	// Get indicators used for strategy

	ss.inputsMtx.RLock()
	sslInput := ss.inputs[indicators.SSLInput]
	ss.inputsMtx.RUnlock()

	ss.indicatorsMtx.RLock()
	sslIndicator := ss.indicators[indicators.SSLIndicatorID]
	ss.indicatorsMtx.RUnlock()

	slog.Info(nil, "Running scalping strategy.")
	for {
		select {
		case e, ok := <-sslInput:
			if !ok {
				ss.Done()
			}
			slog.Info(nil, "Recevied input to SSL Indicator.")
			sslIndicator.Input <- e
		case s := <-sslIndicator.Signals:
			slog.Info(nil, "Signal received: %v", s)
		case <-ss.done:
			slog.Info(nil, "Stopping scaling strategy")
			return
		}
	}
}

func (ss *ScalpingStrategy) Done() {
	// Tell indicator it's time to stop.
	defer ss.cancelFunc()
	select {
	case ss.done <- struct{}{}:
		return
	default:
		return
	}
}
