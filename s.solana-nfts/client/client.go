package client

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/monzo/slog"
	"github.com/opentracing/opentracing-go"

	"swallowtail/libraries/gerrors"
	"swallowtail/libraries/transport"
	"swallowtail/s.solana-nfts/dto"
	"swallowtail/s.solana-nfts/marshaling"
	solananftsproto "swallowtail/s.solana-nfts/proto"
)

var (
	client SolanaNFTVendorClient
)

// SolanaNFTVendorClient ...
type SolanaNFTVendorClient interface {
	// Ping pings all vendors to healthcheck connections
	Ping(ctx context.Context) error
	// GetSolanartPriceStatisticsByCollectionID retrieves price data from solanart for a given collection.
	GetSolanartPriceStatisticsByCollectionID(ctx context.Context, req *dto.GetVendorPriceStatisticsByCollectionIDRequest) (*dto.GetSolanartPriceStatisticsByCollectionIDResponse, error)
	// GetMagicEdenPriceStatisticsByCollectionID price data from magic eden for a given collection.
	GetMagicEdenPriceStatisticsByCollectionID(ctx context.Context, req *dto.GetVendorPriceStatisticsByCollectionIDRequest) (*dto.GetMagicEdenPriceStatisticsByCollectionIDResponse, error)
}

// Init initializes are required vendor clients
func Init(ctx context.Context) error {
	http := transport.NewHTTPClient(10*time.Second, nil)
	c := NewVendorClient(ctx, http)

	if err := c.Ping(ctx); err != nil {
		return gerrors.Augment(err, "failed_to_establish_connection_to_all_vendor_clients", nil)
	}

	client = c

	return nil
}

// GetVendorPriceStatisticsByCollectionID ...
func GetVendorPriceStatisticsByCollectionID(ctx context.Context, vendor solananftsproto.SolanaNFTVendor, req *dto.GetVendorPriceStatisticsByCollectionIDRequest, sorting solananftsproto.SolanaNFTSortDirection, limit int) (*dto.GetVendorPriceStatisticsByCollectionIDResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, fmt.Sprintf("Solana NFT vendor price statistics %s", vendor))
	defer span.Finish()

	errParams := map[string]string{
		"vendor": vendor.String(),
	}

	// Gather data by vendor.
	var rsp *dto.GetVendorPriceStatisticsByCollectionIDResponse
	switch vendor {
	case solananftsproto.SolanaNFTVendor_MAGIC_EDEN:
		meRsp, err := client.GetMagicEdenPriceStatisticsByCollectionID(ctx, req)
		if err != nil {
			return nil, gerrors.Augment(err, "failed_to_get_magic_eden_statistics", errParams)
		}

		rsp = marshaling.MagicEdenPriceStatisticsDTOToVendorDTO(meRsp)
	case solananftsproto.SolanaNFTVendor_SOLANART:
		sRsp, err := client.GetSolanartPriceStatisticsByCollectionID(ctx, req)
		if err != nil {
			return nil, gerrors.Augment(err, "failed_to_get_solanart_statistics", errParams)
		}

		rsp = marshaling.SolanartPriceStatisticsDTOToVendorDTO(sRsp)
	default:
		return nil, gerrors.Unimplemented("failed_to_get_vendor_statistics.unimplemented", errParams)
	}

	if len(rsp.Stats) == 0 {
		slog.Warn(ctx, "Solana NFT collection [%s] from vendor [%s], price stats returned zero data points", req.CollectionID, vendor)
		return rsp, nil
	}

	// Apply sorting.
	switch sorting {
	case solananftsproto.SolanaNFTSortDirection_ASCENDING:
		sort.Slice(rsp.Stats, func(i int, j int) bool {
			return rsp.Stats[i].Price > rsp.Stats[j].Price
		})
	default:
		sort.Slice(rsp.Stats, func(i int, j int) bool {
			return rsp.Stats[i].Price < rsp.Stats[j].Price
		})
	}

	// Apply limit.
	if limit > 0 {
		rsp.Stats = rsp.Stats[0:1]
	}

	return rsp, nil
}
