package execution

import (
	"fmt"
	"sync"

	tradeengineproto "swallowtail/s.trade-engine/proto"
)

var (
	executionRegistry = make(map[tradeengineproto.EXECUTION_STRATEGY]StrategyExecution)
	executionMu       sync.RWMutex
)

func register(strategy tradeengineproto.EXECUTION_STRATEGY, handler StrategyExecution) {
	executionMu.Lock()
	defer executionMu.Unlock()

	if _, ok := executionRegistry[strategy]; ok {
		panic(fmt.Sprintf("Failed to register execution strategy: strategy already registered; %s", strategy))
	}

	executionRegistry[strategy] = handler
}

func getStrategyExecution(strategy tradeengineproto.EXECUTION_STRATEGY) (StrategyExecution, bool) {
	executionMu.Lock()
	defer executionMu.Unlock()

	es, ok := executionRegistry[strategy]
	if !ok {
		return nil, false
	}

	return es, true
}
