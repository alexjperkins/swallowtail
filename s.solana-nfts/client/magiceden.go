package client

import (
	"context"
	"net/http"

	"swallowtail/libraries/gerrors"
	"swallowtail/libraries/ratelimit"
	"swallowtail/libraries/transport"
	"swallowtail/s.solana-nfts/dto"
)

type magicEdenClient struct {
	rateLimiter ratelimit.RateLimiter
	http        transport.HttpClient
}

func (m *magicEdenClient) Ping(ctx context.Context) error {
	// TODO
	//if _, err := m.GetMagicEdenPriceStatisticsByCollectionID(ctx, &dto.GetVendorPriceStatisticsByCollectionIDRequest{
	//	CollectionID: solananftsproto.MagicEndCollectionIDGloomPunks,
	//}); err != nil {
	//	return gerrors.Augment(err, "failed_to_establish_connection_magic_eden", nil)
	//}

	return nil
}

func (m *magicEdenClient) GetMagicEdenPriceStatisticsByCollectionID(
	ctx context.Context, req *dto.GetVendorPriceStatisticsByCollectionIDRequest,
) (*dto.GetMagicEdenPriceStatisticsByCollectionIDResponse, error) {
	endpoint := buildURL(req.CollectionID)

	// Rate limit.
	m.rateLimiter.Throttle()

	rsp := &dto.GetMagicEdenPriceStatisticsByCollectionIDResponse{}
	if err := m.http.Do(ctx, http.MethodGet, endpoint, nil, rsp); err != nil {
		return nil, gerrors.Augment(err, "failed_to_get_price_statistics_magic_eden", nil)
	}

	return nil, nil
}

func buildURL(collectionID string) string {
	return "https://api-mainnet.magiceden.io/rpc/getListedNFTsByQuery?q=%7B%22%24match%22%3A%7B%22collectionSymbol%22%3A%22" + collectionID + "%22%7D%2C%22%24sort%22%3A%7B%22takerAmount%22%3A1%2C%22createdAt%22%3A-1%7D%7D"
}
