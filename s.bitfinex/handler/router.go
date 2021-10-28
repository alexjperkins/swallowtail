package handler

import (
	bitfinexproto "swallowtail/s.bitfinex/proto"
)

// BitfinexService ...
type BitfinexService struct {
	*bitfinexproto.UnimplementedBitfinexServer
}
