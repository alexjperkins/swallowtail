package handler

import (
	tradeengineproto "swallowtail/s.trade-engine/proto"
)

// TradeEngineHandler ...
type TradeEngineHandler struct {
	*tradeengineproto.UnimplementedTradeengineServer
}
