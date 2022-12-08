package handler

import (
	"context"
	"sync"

	"swallowtail/libraries/gerrors"
	"swallowtail/s.market-data/assets"
	marketdataproto "swallowtail/s.market-data/proto"
)

type SolanaNFTInfo struct {
	CollectionID          string
	VendorID              string
	Price                 float64
	HumanizedCollectionID string
	Vendor                string
	Emoji                 string
	Price4H               float64
	Price24H              float64
	TotalListed           int
}

var (
	solanaNFTAssets = assets.SolanaNFTAssets
)

var (
	// solanaNFTPriceCache4h  *ttlcache.Cache
	// solanaNFTPriceCache24h *ttlcache.Cache
	solanaNFTPriceOnce sync.Once
)

// PublishSolanaNFTPriceInformation ...
func (s *MarketDataService) PublishSolanaNFTPriceInformation(
	ctx context.Context, in *marketdataproto.PublishSolanaNFTPriceInformationRequest,
) (*marketdataproto.PublishSolanaNFTPriceInformationResponse, error) {
	return nil, gerrors.ErrUnimplemented
}
