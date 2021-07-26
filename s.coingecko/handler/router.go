package handler

import (
	coingeckoproto "swallowtail/s.coingecko/proto"
)

type CoingeckoService struct {
	*coingeckoproto.UnimplementedCoingeckoServer
}
