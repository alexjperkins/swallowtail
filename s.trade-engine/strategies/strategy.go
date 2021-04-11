package strategies

import (
	"context"
	"time"
)

type Strategy interface {
	Run(context.Context)
	Done()
}

// ParallelStrategy runs strategies in parallel; the order doesn't matter.
type ParallelStrategy struct {
	strategies                   []*Strategy
	confidenceLevelBeforeTrigger int
}

func (ps *ParallelStrategy) Run(context.Context) {}
func (ps *ParallelStrategy) Done()               {}

// SequentialStrategy runs strategies in sequence, it only check the next indicator if the
// prior indeed generates a signal
type SequentialStrategy struct {
	strategies []*Strategy
	// A duration of 0 means wait until told otherwise
	indicatorSignalTimeout time.Duration
}

func (ss *SequentialStrategy) Run(context.Context) {}
func (ss *SequentialStrategy) Done()               {}
