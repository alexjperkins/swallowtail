package handler

import (
	binanceproto "swallowtail/s.binance/proto"
)

type BinanceService struct {
	binanceproto.UnimplementedBinanceServer
}
