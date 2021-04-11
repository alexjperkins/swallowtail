package multiplexing

import (
	"context"
)

type MuliplexFilter func(interface{}) (interface{}, bool)

type MultiplexConsumer struct {
	id       int
	Ch       chan interface{}
	Filter   MuliplexFilter
	Metadata map[string]string
}

func NewMultiplexConsumer(bufSize int, filter MuliplexFilter, metadata map[string]string) *MultiplexConsumer {
	// ID find hashing func
	return &MultiplexConsumer{
		Ch:       make(chan interface{}, bufSize),
		Filter:   filter,
		Metadata: metadata,
	}
}

func (mc *MultiplexConsumer) send(ctx context.Context, e interface{}) {
	go func() {
		te, ok := mc.Filter(e)
		if !ok || te == nil {
			// Early exit
			return
		}
		select {
		case mc.Ch <- te:
			return
		case <-ctx.Done():
			return
		}
	}()
}

func (mc *MultiplexConsumer) close() {
	close(mc.Ch)
}
