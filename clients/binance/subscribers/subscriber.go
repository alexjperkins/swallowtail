package subscribers

import (
	"context"
	"swallowtail/clients/binance/domain"
)

// Subscriber interface
type Subscriber interface {
	Send(context.Context, *domain.BinanceMsg)
	Close()
}
