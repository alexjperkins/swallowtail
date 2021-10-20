package handler

import (
	marketdataproto "swallowtail/s.market-data/proto"
)

// MarketDataService ...
type MarketDataService struct {
	*marketdataproto.UnimplementedMarketdataServer
}
