package handler

import (
	tradeengineproto "swallowtail/s.trade-engine/proto"
)

// TradeEngineService ...
type TradeEngineService struct {
	*tradeengineproto.UnimplementedTradeengineServer
}
