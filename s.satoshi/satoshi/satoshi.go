package satoshi

import "context"

type SatoshiComponents interface {
	Start(ctx context.Context) error
	Done()
	ID() string
}

func New(components ...SatoshiComponents) *Satoshi {
	return &Satoshi{}
}

type Satoshi struct{}
