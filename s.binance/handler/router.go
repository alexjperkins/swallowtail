package handler

import (
	binanceproto "swallowtail/s.binance/proto"
)

// BinanceService defines the service for interacting & making trades with Binance.
type BinanceService struct {
	binanceproto.UnimplementedBinanceServer
}
