package handler

import (
	"context"

	"swallowtail/libraries/gerrors"
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

// PublishSolanaNFTPriceInformation ...
func (s *MarketDataService) PublishSolanaNFTPriceInformation(
	ctx context.Context, in *marketdataproto.PublishSolanaNFTPriceInformationRequest,
) (*marketdataproto.PublishSolanaNFTPriceInformationResponse, error) {
	return nil, gerrors.Unimplemented("publish solana nft price information", nil)
}
