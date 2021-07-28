package handler

import (
	marketdataproto "swallowtail/s.market-data/proto"
)

type MarketDataService struct {
	*marketdataproto.UnimplementedBybtServer
}
