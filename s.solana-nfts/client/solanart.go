package client

import (
	"context"
	"fmt"
	"net/http"

	"swallowtail/libraries/gerrors"
	"swallowtail/libraries/ratelimit"
	"swallowtail/libraries/transport"
	"swallowtail/s.solana-nfts/dto"
	solananftsproto "swallowtail/s.solana-nfts/proto"
)

const (
	solanartURL = "https://qzlsklfacc.medianetwork.cloud"
)

type solanartClient struct {
	rateLimiter ratelimit.RateLimiter
	http        transport.HttpClient
}

func (s *solanartClient) Ping(ctx context.Context) error {
	// Here since the API isn't **public** we just ping a known endpoint.
	if _, err := s.GetSolanartPriceStatisticsByCollectionID(ctx, &dto.GetVendorPriceStatisticsByCollectionIDRequest{
		CollectionID: solananftsproto.SolanartCollectionIDGalacticGeckoSpaceGarage,
	}); err != nil {
		return err
	}

	return nil
}

func (s *solanartClient) GetSolanartPriceStatisticsByCollectionID(
	ctx context.Context, req *dto.GetVendorPriceStatisticsByCollectionIDRequest,
) (*dto.GetSolanartPriceStatisticsByCollectionIDResponse, error) {
	endpoint := fmt.Sprintf("%s/nft_for_sale?collection=%s", solanartURL, req.CollectionID)

	rsp := &dto.GetSolanartPriceStatisticsByCollectionIDResponse{}
	if err := s.http.Do(ctx, http.MethodGet, endpoint, nil, rsp); err != nil {
		return nil, gerrors.Augment(err, "failed_to_get_solanart_price_statistics", nil)
	}

	return rsp, nil
}
